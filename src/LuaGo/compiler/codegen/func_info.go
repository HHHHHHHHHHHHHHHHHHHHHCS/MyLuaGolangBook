package codegen

import (
	. "LuaGo/compiler/ast"
	. "LuaGo/vm"
)

type locVarInfo struct {
	prev    *locVarInfo //上一层
	name    string      //局部变量名字
	scopeLv int         //作用域层级
	slot    int         //绑定的索引
	capture bool        //局部变量是否引用的父亲(外部)
}

type funcInfo struct {
	constants map[interface{}]int    //常量表
	usedRegs  int                    //已经使用的寄存器
	maxRegs   int                    //最大寄存器
	scopeLv   int                    //作用域层级
	locVars   []*locVarInfo          //内部申明的全部局部变量
	locNames  map[string]*locVarInfo //当前生效的局部变量
	breaks    [][]int                //跳出块
	parent    *funcInfo              //父func
	upvalues  map[string]upvalInfo   //upValue
	insts     []uint32               //指令
	subFuncs  []*funcInfo            //子func信息 可能有多个
	numParams int                    //参数数量
	isVararg  bool                   //可变参数
}

//upvaltable  唯一的
type upvalInfo struct {
	locVarSlot int
	upvalIndex int
	index      int
}

func newFuncInfo(parent *funcInfo, fd *FuncDefExp) *funcInfo {
	return &funcInfo{
		parent:    parent,
		subFuncs:  []*funcInfo{},
		constants: map[interface{}]int{},
		upvalues:  map[string]upvalInfo{},
		locNames:  map[string]*locVarInfo{},
		locVars:   make([]*locVarInfo, 0, 8),
		breaks:    make([][]int, 1),
		insts:     make([]uint32, 0, 8),
		isVararg:  fd.IsVararg,
		numParams: len(fd.ParList),
	}
}

func (self *funcInfo) indexOfConstant(k interface{}) int {
	if idx, found := self.constants[k]; found {
		return idx
	}

	idx := len(self.constants)
	self.constants[k] = idx
	return idx
}

//分配寄存器
func (self *funcInfo) allocReg() int {
	self.usedRegs++
	//不能超过255 是因为指令的长度限制
	if self.usedRegs > 255 {
		panic("function or expression needs to many registers")
	}
	if self.usedRegs > self.maxRegs {
		self.maxRegs = self.usedRegs
	}
	return self.usedRegs - 1
}

//分配n个
func (self *funcInfo) allocRegs(n int) int {
	for i := 0; i < n; i++ {
		self.allocReg()
	}
	return self.usedRegs - n
}

//回收一个
func (self *funcInfo) freeReg() {
	self.usedRegs--
}

//回收n个
func (self *funcInfo) freeRegs(n int) {
	for i := 0; i < n; i++ {
		self.freeReg()
	}
}

//进入进一步作用域
func (self *funcInfo) enterScope(breakable bool) {
	self.scopeLv++
	if breakable {
		self.breaks = append(self.breaks, []int{}) //循环块
	} else {
		self.breaks = append(self.breaks, nil) //非循环块
	}
}

//添加一个局部变量
func (self *funcInfo) addLocVar(name string) int {
	newVar := &locVarInfo{
		name:    name,
		prev:    self.locNames[name],
		scopeLv: self.scopeLv,
		slot:    self.allocReg(),
	}
	self.locVars = append(self.locVars, newVar)
	self.locNames[name] = newVar
	return newVar.slot
}

//是否存在这个局部变量 没有返回-1
func (self *funcInfo) slotOfLocVar(name string) int {
	if locVar, found := self.locNames[name]; found {
		return locVar.slot
	}
	return -1
}

func (self *funcInfo) removeLocVar(locVar *locVarInfo) {
	self.freeReg()
	if locVar.prev == nil {
		delete(self.locNames, locVar.name)
	} else if locVar.prev.scopeLv == locVar.scopeLv {
		self.removeLocVar(locVar.prev) //递归删除
	} else {
		self.locNames[locVar.name] = locVar.prev
	}
}

func (self *funcInfo) exitScope() {
	pendingBreakJmps := self.breaks[len(self.breaks)-1]
	self.breaks = self.breaks[:len(self.breaks)-1]
	a := self.getJmpArgA()
	for _, pc := range pendingBreakJmps {
		sBx := self.pc() - pc
		i := (sBx+MAXARG_sBx)<<14 | a<<6 | OP_JMP
		self.insts[pc] = uint32(i)
	}

	self.scopeLv--
	for _, locVar := range self.locNames {
		if locVar.scopeLv > self.scopeLv { //离开作用域
			self.removeLocVar(locVar)
		}
	}
}

//添加break跳转
func (self *funcInfo) addBreakJmp(pc int) {
	for i := self.scopeLv; i >= 0; i-- {
		if self.breaks[i] != nil {
			self.breaks[i] = append(self.breaks[i], pc)
			return
		}
	}
	panic("<break> at line ? not inside a loop!")
}

//寻找upvaltable
func (self *funcInfo) indexOfUpval(name string) int {
	if upval, ok := self.upvalues[name]; ok {
		return upval.index
	}

	if self.parent != nil {
		if locVar, found := self.parent.locNames[name]; found {
			idx := len(self.upvalues)
			self.upvalues[name] = upvalInfo{locVar.slot, -1, idx}
			locVar.capture = true
			return idx
		}

		if uvIdx := self.parent.indexOfUpval(name); uvIdx >= 0 {
			idx := len(self.upvalues)
			self.upvalues[name] = upvalInfo{-1, uvIdx, idx}
			return idx
		}
	}
	return -1
}



func (self *funcInfo) getJmpArgA() int {
	hasCapturedLocVars := false
	minSlotOffLocVars := self.maxRegs
	for _, locVar := range self.locNames {
		if locVar.scopeLv == self.scopeLv {
			for v := locVar; v != nil && v.scopeLv == self.scopeLv; v = v.prev {
				if v.capture {
					hasCapturedLocVars = true
				}
				if v.slot < minSlotOffLocVars && v.name[0] != '(' {
					minSlotOffLocVars = v.slot
				}
			}
		}
	}

	if hasCapturedLocVars {
		return minSlotOffLocVars + 1
	} else {
		return 0
	}
}

func (self *funcInfo) emitABC(opcode, a, b, c int) {
	i := b<<23 | c<<14 | a<<6 | opcode
	self.insts = append(self.insts, uint32(i))
}

func (self *funcInfo) emitABx(opcode, a, bx int) {
	i := bx<<14 | a<<6 | opcode
	self.insts = append(self.insts, uint32(i))
}

func (self *funcInfo) emitAsBx(opcode, a, b int) {
	i := (b+MAXARG_sBx)<<14 | a<<6 | opcode
	self.insts = append(self.insts, uint32(i))
}

func (self *funcInfo) emitAx(opcode, ax int) {
	i := ax<<6 | opcode
	self.insts = append(self.insts, uint32(i))
}

func (self *funcInfo) pc() int {
	return len(self.insts) - 1
}

func (self *funcInfo) fixSbx(pc, sBx int) {
	i := self.insts[pc]
	i = i << 18 >> 18                  //清除sbx操作数
	i = i | uint32(sBx+MAXARG_sBx)<<14 //重置sbx操作数
	self.insts[pc] = 1
}
