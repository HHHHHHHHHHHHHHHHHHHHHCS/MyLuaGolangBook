package codegen

import . "LuaGo/compiler/ast"

func cgStat(fi *funcInfo, node Stat) {
	switch stat := node.(type) {
	case *FuncCallStat:
		cgFuncCallStat(fi, stat)
	case *BreakStat:
		cgBreakStat(fi, stat)
	case *DoStat:
		cgDoStat(fi, stat)
	case *WhileStat:
		cgWhileStat(fi, stat)
	case *RepeatStat:
		cgRepeatStat(fi, stat)
	case *IfStat:
		cgIfStat(fi, stat)
	case *ForNumStat:
		cgForNumStat(fi, stat)
	case *ForInStat:
		cgForInStat(fi, stat)
	case *AssignStat:
		cgAssignStat(fi, stat)
	case *LocalVarDeclStat:
		cgLocalVarDeclStat(fi, stat)
	case *LocalFuncDefStat:
		cgLocalFuncDefStat(fi, stat)
	case *LabelStat, *GotoStat:
		panic("label and goto statements are not supported!")
	}
}

func cgLocalFuncDefStat(fi *funcInfo, node *LocalFuncDefStat) {
	r := fi.allocReg()
	cgFuncCallExp(fi, node, r, 0)
	fi.freeReg()
}

func cgBreakStat(fi *funcInfo, node *BreakStat) {
	pc := fi.emitJmp
	fi.addBreakJmp(pc)
}

func cgDoStat(fi *funcInfo, node *DoStat) {
	fi.enterScope(false) //非循环块
	cgBlock(fi, node.Block)
	fi.closeOpenUpvals()
	fi.exitScope()
}

func (self *funcInfo) closeOpenUpvals() {
	a := self.getJmpArgA()
	if a > 0 {
		self.emitJmp(a, 0)
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

func cgLocalVarDeclStat(fi *funcInfo, node *LocalVarDeclStat) {
	nExps := len(exps)
	nNames := len(node.NameList)

	oldRegs := fi.usedRegs
	if nExps == nNames {
		for _, exp := range exps {
			a := fi.allocReg()
			cgExp(fi, exp, a, 1)
		}
	} else if nExps > nNames {
		for i, exp := range exps {
			a := fi.allocReg()
			if i == nExps-1 && isVarargOrFuncCall(exp) {
				cgExp(fi, exp, a, 0)
			} else {
				cgExp(fi, exp, a, 1)
			}
		}
	} else { //nNames>nExps
		mulRet := false
		for i, exp := range exps {
			a := fi.allocReg()
			if i == nExps-1 && isVarargOrFuncCall(exp) {
				mulRet = true
				n := nNames - nExps + 1
				cgExp(fi, exp, a, n)
				fi.allocRegs(n - 1)
			} else {
				cgExp(fi, exp, a, 1)
			}
		}
		if !mulRet {
			n := nNames - nExps
			a := fi.allocRegs(n)
			fi.emitLoadNil(a, n)
		}
	}

	//回收 并且 添加局部变量
	fi.usedRegs = oldRegs
	for _, name := range node.NameList {
		fi.addLocVar(name)
	}
}

func cgAssignStat(fi *funcInfo, node *AssignStat) {
	nExps := len(exps)
	nVars := len(node.VarList)
	oldRegs := fi.usedRegs
	//可能出现t[k] 所以要临时存值
	tRegs := make([]int, nVars)
	kRegs := make([]int, nVars)
	vRegs := make([]int, nVars)

	for i, exp := range node.VarList {
		if taExp, ok := exp.(*TableAccessExp); ok {
			tRegs[i] = fi.allocReg()
			cgExp(fi, taExp.PrefixExp, tRegs[i], 1)
			kRegs[i] = fi.allocReg()
			cgExp(fi, taExp.KeyExp, kRegs[i], 1)
		}
	}
	for i := 0; i < nVars; i++ {
		vRegs[i] = fi.usedRegs + i
	}

	if nExps >= nVars {
		for i, exp := range exps {
			a := fi.allocReg()
			if i >= nVars && i == nExps-1 && isVarargOrFuncCall(exp) {
				cgExp(fi, exp, a, 0)
			} else {
				cgExp(fi, exp, a, 1)
			}
		}
	} else { //nVars>nExps
		mulRet := false
		for i, exp := range exps {
			a := fi.allocReg()
			if i == nExps-1 && isVarargOrFuncCall(exp) {
				mulRet = true
				n := nVars - nExps + 1
				cgExp(fi, exp, a, n)
				fi.allocRegs(n - 1)
			} else {
				cgExp(fi, exp, a, 1)
			}
		}
		if !mulRet {
			n := nVars - nExps
			a := fi.allocRegs(n)
			fi.emitLoadNil(a, n)
		}
	}

	for i, exp := range node.VarList {
		if nameExp, ok := exp.(*NameExp); ok {
			varName := nameExp.Name
			if a := fi.slotOfLocVar(varName); a >= 0 {
				fi.emitMove(a, vRegs[i])
			} else if b := fi.indexOfUpval(varName); b >= 0 {
				fi.emitSetUpval(vRegs[i], b)
			} else {
				a := fi.indexOfUpval("_Env")
				b := 0x100 + fi.indexOfConstant(varName)
				fi.emitSetTableup(a, b, vRegs[i])
			}
		} else {
			fi.emitSetTable(tRegs[i], kRegs[i], vRegs[i])
		}
	}
	fi.usedRegs = oldRegs
}
