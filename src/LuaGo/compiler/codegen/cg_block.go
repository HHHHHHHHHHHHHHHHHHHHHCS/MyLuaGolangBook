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

//	exp->true->block ->while循环
//		->false->other
func cgWhileStat(fi *funcInfo, node *WhileStat) {
	pcBeforeExp := fi.pc()
	//第2步
	r := fi.allocReg()
	cgExp(fi, node.Exp, r, 1)
	fi.freeReg()
	//第3步
	fi.emitTest(r, 0)
	pcJmpToEnd := fi.emitJmp(0, 0)
	//第4步
	fi.enterScope(true)
	cgBlock(fi, node.Block)
	fi.closeOpenUpvals()
	fi.emitJmp(0, pcBeforeExp-fi.pc()-1)
	fi.exitScope()
	//第5步
	fi.fixSbx(pcJmpToEnd, fi.pc()-pcJmpToEnd)
}

//	block->exp->true->
//				->false->repeat循环
func cgRepeatStat(fi *funcInfo, node *RepeatStat) {
	fi.enterScope(true)

	pcBeforeBlock := fi.pc()
	cgBlock(fi, node.Block)

	r := fi.allocReg()
	cgExp(fi, node.Exp, r, 1)
	fi.freeReg()

	fi.emitTest(r, 0)
	fi.emitJmp(fi.getJmpArgA(), pcBeforeBlock-fi.pc()-1)
	fi.closeOpenUpvals()

	fi.exitScope()
}

func cgIfStat(fi *funcInfo, node *IfStat) {
	pcJmpToEnds := make([]int, len(node.Exps))
	pcJmpToNextExp := -1

	for i, exp := range node.Exps {
		if pcJmpToNextExp >= 0 {
			fi.fixSbx(pcJmpToNextExp, fi.pc()-pcJmpToNextExp)
		}

		r := fi.allocReg()
		cgExp(fi, exp, r, 1)
		fi.freeReg()

		fi.emitTest(r, 0)
		pcJmpToNextExp = fi.emitJmp(0, 0)

		fi.enterScope(false)
		cgBlock(fi, node.Blocks[i])
		fi.closeOpenUpvals()
		fi.exitScope()

		if i < len(node.Exps)-1 {
			pcJmpToEnds[i] = fi.emitJmp(0, 0)
		} else {
			pcJmpToEnds[i] = pcJmpToNextExp
		}
	}

	for _, pc := range pcJmpToEnds {
		fi.fixSbx(pc, fi.pc()-pc)
	}
}

func cgForNumStat(fi *funcInfo, node ForNumStat) {
	fi.enterScope(true)

	//1.
	cgLocalVarDeclStat(fi, &LocalVarDeclStat{
		NameList: []string{"(for index)", "(for limit)", "(for step)"},
		ExpList:  []Exp{node.InitExp, node.LimitExp, node.StepExp},
	})
	fi.addLocVar(node.VarName)

	//2.
	a := fi.usedRegs - 4
	pcForPrep := fi.emitForPrep(a, 0)
	cgBlock(fi, node.Block)
	fi.closeOpenUpvals()
	pcForLoop := fi.emitForLoop(a, 0)

	//3.
	fi.fixSbx(pcForPrep, pcForLoop-pcForPrep-1)

	fi.exitScope()
}

func cgForInStat(fi *funcInfo, node *ForInStat) {
	fi.enterScope(true)
	//1.
	cgLocalVarDeclStat(fi, &LocalVarDeclStat{
		NameList: []string{"(for generator)", "(for state)", "(for control)"},
		ExpList:  node.ExpList,
	})
	for _, name := range node.NameList {
		fi.addLocVar(name)
	}

	//2.
	pcJmpToTFC := fi.emitJmp(0, 0)
	cgBlock(fi, node.Block)
	fi.closeOpenUpvals()
	fi.fixSbx(pcJmpToTFC, fi.pc()-pcJmpToTFC)

	//3.
	rGenerator := fi.slotOfLocVar("(for generator)")
	fi.emitTForCall(rGenerator, len(node.NameList))
	fi.emitTForLoop(rGenerator+2, pcJmpToTFC-fi.pc()-1)

	fi.exitScope()
}
