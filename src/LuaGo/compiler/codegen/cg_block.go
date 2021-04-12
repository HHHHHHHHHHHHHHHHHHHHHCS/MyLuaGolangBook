package codegen

import . "LuaGo/compiler/ast"

func cgBlock(fi *funcInfo, node *Block) {
	for _, stat := range node.Stats {
		cgStat(fi, stat)
	}

	if node.RetExps != nil {
		cgRetStat(fi, node.RetExps)
	}
}

func isVarargOrFuncCall(exp Exp) bool {
	switch exp.(type) {
	case *VarargExp, *FuncCallExp:
		return true
	}
	return false
}

func cgRetStat(fi *funcInfo, exps []Exp) {
	nExps := len(exps)
	if nExps == 0 { //直接使用return 后面没有东西了
		fi.emitReturn(0, 0)
		return
	}

	multRet := isVarargOrFuncCall(exps[nExps-1])
	for i, exp := range exps {
		r := fi.allocReg()
		//TODO:没有处理尾递归调用
		if i == nExps-1 && multRet {
			//多个返回值
			cgExp(fi, exp, r, -1)
		} else {
			//单个返回值
			cgExp(fi, exp, r, 1)
		}
	}

	//回收
	fi.freeRegs(nExps)
	a := fi.usedRegs
	if multRet {
		fi.emitReturn(a, -1)
	} else {
		fi.emitReturn(a, nExps)
	}
}


