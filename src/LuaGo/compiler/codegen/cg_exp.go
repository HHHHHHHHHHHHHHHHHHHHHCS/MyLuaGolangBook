package codegen

import . "LuaGo/compiler/ast"

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
