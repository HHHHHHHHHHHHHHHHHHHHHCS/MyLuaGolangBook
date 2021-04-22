package codegen

import (
	. "LuaGo/compiler/ast"
	. "LuaGo/compiler/lexer"
	. "LuaGo/vm"
)

func cgExp(fi *funcInfo, node Exp, a, n int) {
	switch exp := node.(type) {
	case *NilExp:
		fi.emitLoadNil(a, n)
		break
	case *FalseExp:
		fi.emitLoadBool(a, 0, 0)
		break
	case *TrueExp:
		fi.emitLoadBool(a, 1, 0)
		break
	case *IntegerExp:
		fi.emitLoadK(a, exp.Val)
		break
	case *FloatExp:
		fi.emitLoadK(a, exp.Val)
		break
	case *StringExp:
		fi.emitLoadK(a, exp.str)
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
		cgNameExp(fi, exp, a)
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
	fi.emitVararg(a, n)
}

func cgFuncDefExp(fi *funcInfo, node *FuncDefExp, a int) {
	subFI := newFuncInfo(fi, node)
	fi.subFuncs = append(fi.subFuncs, subFI)

	for _, param := range node.ParList {
		subFI.addLocVar(param)
	}

	cgBlock(subFI, node.Block)
	subFI.exitScope()
	subFI.emitReturn(0, 0) //lua给每一个函数都添加了return

	bx := len(fi.subFuncs) - 1
	fi.emitClosure(a, bx) //退出
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

	fi.emitNewTable(a, nArr, nExps-nArr)

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
				c := (arrIdx-1)/50 + 1
				fi.freeRegs(n)
				if i == nExps-1 && multRet {
					fi.emitSetList(a, 0, c)
				} else {
					fi.emitSetList(a, n, c)
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
		fi.emitSetTable(a, b, c)
	}
}

//一元表达式
func cgUnopExp(fi *funcInfo, node *UnopExp, a int) {
	b := fi.allocReg()
	cgExp(fi, node.Exp, b, 1)
	fi.emitUnaryOp(node.Op, a, b)
	fi.freeReg()
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
	fi.emitABC(OP_CONCAT, a, b, c)
}

//逻辑表达式  操作数需要特殊对待 需要生成testset 和 move
//其它的生成临时变量求值就好了
func cgBinopExp(fi *funcInfo, node *BinopExp, a int) {
	switch node.Op {
	case TOKEN_OP_ADD, TOKEN_OP_OR:
		b := fi.allocReg()
		cgExp(fi, node.Exp1, b, 1)
		fi.freeReg()
		if node.Op == TOKEN_OP_AND {
			fi.emitTestSet(a, b, 0)
		} else {
			fi.emitTestSet(a, b, 1)
		}
		pcOfJmp := fi.emitJmp(0, 0)

		b = fi.allocReg()
		cgExp(fi, node.Exp2, b, 1)
		fi.freeReg()
		fi.emitMove(a, b)
		fi.fixSbx(pcOfJmp, fi.pc()-pcOfJmp)
	default:
		b := fi.allocReg()
		cgExp(fi, node.Exp1, b, 1)
		c := fi.allocReg()
		cgExp(fi, node.Exp2, c, 1)
		fi.emitBinaryOp(node.Op, a, b, c)
		fi.freeRegs(2)
	}
}

func cgTableAccessExp(fi *funcInfo, node *TableAccessExp, a int) {
	b := fi.allocReg()
	cgExp(fi, node.PrefixExp, b, 1)
	c := fi.allocReg()
	cgExp(fi, node.KeyExp, c, 1)
	fi.emitGetTable(a, b, c)
	fi.freeRegs(2)
}

//局部变量则move     upvalue则getupval    全局变量则表访问
func cgNameExp(fi *funcInfo, node *NameExp, a int) {
	if r := fi.slotOfLocVar(node.Name); r >= 0 {
		fi.emitMove(a, r)
	} else if idx := fi.indexOfUpval(node.Name); idx >= 0 {
		fi.emitGetUpval(a, idx)
	} else { //x=>_ENV['x']
		taExp := &TableAccessExp{
			PrefixExp: &NameExp{Name: "_Env"},
			KeyExp:    &StringExp{Str: node.Name},
		}
		cgTableAccessExp(fi, taExp, a)
	}
}

//funcall 获取n个参数 call
func cgFuncCallExp(fi *funcInfo, node *FuncCallExp, a, n int) {
	nArgs := prepFuncCall(fi, node, a)
	fi.emitCall(a, nArgs, n)
}


func prepFuncCall(fi *funcInfo, node *FuncCallExp, a int) int {
	nArgs := len(node.Args)
	lastArgIsVarargOrFuncCall := false

	cgExp(fi, node.PrefixExp, a, 1)

	if node.NameExp != nil {
		c := 0x100 + fi.indexOfConstant(node.NameExp.Str)
		fi.emitSelf(a, a, c)
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
		nArgs++
	}

	if lastArgIsVarargOrFuncCall {
		nArgs = -1
	}

	return nArgs
}

