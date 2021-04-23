package codegen

import (
	. "LuaGo/compiler/ast"
	. "LuaGo/compiler/lexer"
	. "LuaGo/vm"
)

// kind of operands
const (
	ARG_CONST = 1 // const index
	ARG_REG   = 2 // register index
	ARG_UPVAL = 4 // upvalue index
	ARG_RK    = ARG_REG | ARG_CONST
	ARG_RU    = ARG_REG | ARG_UPVAL
	ARG_RUK   = ARG_REG | ARG_UPVAL | ARG_CONST
)

func cgExp(fi *funcInfo, node Exp, a, n int) {
	switch exp := node.(type) {
	case *NilExp:
		fi.emitLoadNil(exp.Line, a, n)
		break
	case *FalseExp:
		fi.emitLoadBool(exp.Line, a, 0, 0)
		break
	case *TrueExp:
		fi.emitLoadBool(exp.Line, a, 1, 0)
		break
	case *IntegerExp:
		fi.emitLoadK(exp.Line, a, exp.Val)
		break
	case *FloatExp:
		fi.emitLoadK(exp.Line, a, exp.Val)
		break
	case *StringExp:
		fi.emitLoadK(exp.Line, a, exp.Str)
		break
	case *ParensExp:
		cgExp(fi, exp.Exp, a, 1)
		break
	case *VarargExp:
		cgVarargExp(fi, exp, a, n)
		break
	case *FuncDefExp:
		cgFuncDefExp(fi, exp, a)
		break
	case *TableConstructorExp:
		cgTableConstructorExp(fi, exp, a)
		break
	case *UnopExp:
		cgUnopExp(fi, exp, a)
		break
	case *BinopExp:
		cgBinopExp(fi, exp, a)
		break
	case *ConcatExp:
		cgConcatExp(fi, exp, a)
		break
	case *NameExp:
		cgNameExp(fi, exp, a)
		break
	case *TableAccessExp:
		cgTableAccessExp(fi, exp, a)
		break
	case *FuncCallExp:
		cgFuncCallExp(fi, exp, a, n)
		break
	}
}

func cgVarargExp(fi *funcInfo, node *VarargExp, a, n int) {
	if !fi.isVararg {
		panic("cannot use '...' outside a vararg function")
	}
	fi.emitVararg(node.Line, a, n)
}

func cgFuncDefExp(fi *funcInfo, node *FuncDefExp, a int) {
	subFI := newFuncInfo(fi, node)
	fi.subFuncs = append(fi.subFuncs, subFI)

	for _, param := range node.ParList {
		subFI.addLocVar(param, 0)
	}

	cgBlock(subFI, node.Block)
	subFI.exitScope(subFI.pc() + 2)
	subFI.emitReturn(node.LastLine, 0, 0) //lua给每一个函数都添加了return

	bx := len(fi.subFuncs) - 1
	fi.emitClosure(node.LastLine, a, bx) //退出
}

//表构造表达式
func cgTableConstructorExp(fi *funcInfo, node *TableConstructorExp, a int) {
	nArr := 0
	for _, keyExp := range node.KeyExps {
		if keyExp == nil {
			nArr++
		}
	}
	nExps := len(node.KeyExps)
	multRet := nExps > 0 && isVarargOrFuncCall(node.ValExps[nExps-1])

	fi.emitNewTable(node.Line, a, nArr, nExps-nArr)

	arrIdx := 0

	for i, keyExp := range node.KeyExps {
		valExp := node.ValExps[i]

		if keyExp == nil { //如果没有key则for 循环增加数组
			arrIdx++
			tmp := fi.allocReg()
			if i == nExps-1 && multRet {
				cgExp(fi, valExp, tmp, -1)
			} else {
				cgExp(fi, valExp, tmp, 1)
			}

			if arrIdx%50 == 0 || arrIdx == nArr { //LFIELDS_PER_FLUSH
				n := arrIdx % 50
				if n == 0 {
					n = 50
				}
				fi.freeRegs(n)
				line := lastLineOf(valExp)
				c := (arrIdx-1)/50 + 1 //todo: c>0xFF
				if i == nExps-1 && multRet {
					fi.emitSetList(line, a, 0, c)
				} else {
					fi.emitSetList(line, a, n, c)
				}
			}

			continue
		}

		//关联表
		b := fi.allocReg()
		cgExp(fi, keyExp, b, 1)
		c := fi.allocReg()
		cgExp(fi, valExp, c, 1)
		fi.freeRegs(2)

		line := lastLineOf(valExp)
		fi.emitSetTable(line, a, b, c)
	}
}

//一元表达式
func cgUnopExp(fi *funcInfo, node *UnopExp, a int) {
	oldRegs := fi.usedRegs
	b, _ := expToOpArg(fi, node.Exp, ARG_REG)
	fi.emitUnaryOp(node.Line, node.Op, a, b)
	fi.usedRegs = oldRegs
}

//逻辑表达式  操作数需要特殊对待 需要生成testset 和 move
//其它的生成临时变量求值就好了
func cgBinopExp(fi *funcInfo, node *BinopExp, a int) {
	switch node.Op {
	case TOKEN_OP_AND, TOKEN_OP_OR:
		oldRegs := fi.usedRegs

		b, _ := expToOpArg(fi, node.Exp1, ARG_REG)
		fi.usedRegs = oldRegs

		if node.Op == TOKEN_OP_AND {
			fi.emitTestSet(node.Line, a, b, 0)
		} else {
			fi.emitTestSet(node.Line, a, b, 1)
		}
		pcOfJmp := fi.emitJmp(node.Line, 0, 0)

		b, _ = expToOpArg(fi, node.Exp2, ARG_REG)
		fi.usedRegs = oldRegs
		fi.emitMove(node.Line, a, b)
		fi.fixSbx(pcOfJmp, fi.pc()-pcOfJmp)
	default:
		oldRegs := fi.usedRegs
		b, _ := expToOpArg(fi, node.Exp1, ARG_RK)
		c, _ := expToOpArg(fi, node.Exp2, ARG_RK)
		fi.emitBinaryOp(node.Line, node.Op, a, b, c)
		fi.usedRegs = oldRegs
	}
}

//拼接字符串
func cgConcatExp(fi *funcInfo, node *ConcatExp, a int) {
	for _, subExp := range node.Exps {
		a := fi.allocReg()
		cgExp(fi, subExp, a, 1)
	}

	c := fi.usedRegs - 1
	b := c - len(node.Exps) + 1
	fi.freeRegs(c - b + 1)
	fi.emitABC(node.Line, OP_CONCAT, a, b, c)
}

//局部变量则move     upvalue则getupval    全局变量则表访问
func cgNameExp(fi *funcInfo, node *NameExp, a int) {
	if r := fi.slotOfLocVar(node.Name); r >= 0 {
		fi.emitMove(node.Line, a, r)
	} else if idx := fi.indexOfUpval(node.Name); idx >= 0 {
		fi.emitGetUpval(node.Line, a, idx)
	} else { //x=>_ENV['x']
		taExp := &TableAccessExp{
			LastLine:  node.Line,
			PrefixExp: &NameExp{Line: node.Line, Name: "_ENV"},
			KeyExp:    &StringExp{Line: node.Line, Str: node.Name},
		}
		cgTableAccessExp(fi, taExp, a)
	}
}

//table
func cgTableAccessExp(fi *funcInfo, node *TableAccessExp, a int) {
	oldRegs := fi.usedRegs
	b, kindB := expToOpArg(fi, node.PrefixExp, ARG_RU)
	c, _ := expToOpArg(fi, node.KeyExp, ARG_RK)
	fi.usedRegs = oldRegs

	if kindB == ARG_UPVAL {
		fi.emitGetTabUp(node.LastLine, a, b, c)
	} else {
		fi.emitGetTable(node.LastLine, a, b, c)

	}
}

//funcall 获取n个参数 call
func cgFuncCallExp(fi *funcInfo, node *FuncCallExp, a, n int) {
	nArgs := prepFuncCall(fi, node, a)
	fi.emitCall(node.Line, a, nArgs, n)
}

//return f(args)
func cgTailCallExp(fi *funcInfo, node *FuncCallExp, a int) {
	nArgs := prepFuncCall(fi, node, a)
	fi.emitTailCall(node.Line, a, nArgs)
}

func prepFuncCall(fi *funcInfo, node *FuncCallExp, a int) int {
	nArgs := len(node.Args)
	lastArgIsVarargOrFuncCall := false

	cgExp(fi, node.PrefixExp, a, 1)

	if node.NameExp != nil {
		fi.allocReg()
		c, k := expToOpArg(fi, node.NameExp, ARG_RK)
		fi.emitSelf(node.Line, a, a, c)
		if k == ARG_REG {
			fi.freeRegs(1)
		}
	}

	for i, arg := range node.Args {
		tmp := fi.allocReg()
		if i == nArgs-1 && isVarargOrFuncCall(arg) {
			lastArgIsVarargOrFuncCall = true
			cgExp(fi, arg, tmp, -1)
		} else {
			cgExp(fi, arg, tmp, 1)
		}
	}

	fi.freeRegs(nArgs)

	if node.NameExp != nil {
		fi.freeReg()
		nArgs++
	}

	if lastArgIsVarargOrFuncCall {
		nArgs = -1
	}

	return nArgs
}

func expToOpArg(fi *funcInfo, node Exp, argKinds int) (arg, argKind int) {
	if argKinds&ARG_CONST > 0 {
		idx := -1
		switch x := node.(type) {
		case *NilExp:
			idx = fi.indexOfConstant(nil)
			break
		case *FalseExp:
			idx = fi.indexOfConstant(false)
			break
		case *TrueExp:
			idx = fi.indexOfConstant(true)
			break
		case *IntegerExp:
			idx = fi.indexOfConstant(x.Val)
			break
		case *FloatExp:
			idx = fi.indexOfConstant(x.Val)
			break
		case *StringExp:
			idx = fi.indexOfConstant(x.Str)
			break
		}

		if idx >= 0 && idx <= 0xFF {
			return 0x100 + idx, ARG_CONST
		}
	}

	if nameExp, ok := node.(*NameExp); ok {
		if argKinds&ARG_REG > 0 {
			if r := fi.slotOfLocVar(nameExp.Name); r >= 0 {
				return r, ARG_REG
			}
		}
		if argKinds&ARG_UPVAL > 0 {
			if idx := fi.indexOfUpval(nameExp.Name); idx >= 0 {
				return idx, ARG_UPVAL
			}
		}
	}

	a := fi.allocReg()
	cgExp(fi, node, a, 1)
	return a, ARG_REG
}
