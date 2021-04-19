package codegen

import (
	. "LuaGo/compiler/ast"
	. "LuaGo/vm"
)

func cgExp(fi *funcInfo, node Exp, a, n int) {
	switch exp := node.(type) {
	case NilExp:
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

func (self *funcInfo) emitNewTable(a, nArr, nRec int) {
	self.emitABC(OP_NEWTABLE, a, Int2fb(nArr), Int2fb(nRec))
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
