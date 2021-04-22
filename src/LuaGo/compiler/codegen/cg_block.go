package codegen

import . "LuaGo/compiler/ast"

func cgBlock(fi *funcInfo, node *Block) {
	for _, stat := range node.Stats {
		cgStat(fi, stat)
	}

	if node.RetExps != nil {
		cgRetStat(fi, node.RetExps, node.LastLine)
	}
}

func cgRetStat(fi *funcInfo, exps []Exp, lastLine int) {
	nExps := len(exps)
	if nExps == 0 { //直接使用return 后面没有东西了
		fi.emitReturn(lastLine, 0, 0)
		return
	}

	if nExps == 1 {
		if nameExp, ok := exps[0].(*NameExp); ok {
			if r := fi.slotOfLocVar(nameExp.Name); r >= 0 {
				fi.emitReturn(lastLine, r, 1)
				return
			}
		}
		if fcExp, ok := exps[0].(*FuncCallExp); ok {
			r := fi.allocReg()
			cgTailCallExp(fi, fcExp, r)
			fi.freeReg()
			fi.emitReturn(lastLine, r, -1)
			return
		}
	}

	multRet := isVarargOrFuncCall(exps[nExps-1])
	for i, exp := range exps {
		r := fi.allocReg()
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
		fi.emitReturn(lastLine, a, -1)
	} else {
		fi.emitReturn(lastLine, a, nExps)
	}
}
