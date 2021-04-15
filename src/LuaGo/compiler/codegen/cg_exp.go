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
	case *ParensSExp:
		cgExp(fi, exp.Exp, a, 1)
		break
	case *VarargExp:
		cgVarargExp(fi, exp, a, n)
		break
	case *FuncDefExp:
		cgFuncDefExp(fi,exp,a)
		break
	case *TableConstructorExp:
		cgTableConstructorExp(fi,exp,a)
		break
	case *UnopExp:
		//todo:
	}
}
