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
