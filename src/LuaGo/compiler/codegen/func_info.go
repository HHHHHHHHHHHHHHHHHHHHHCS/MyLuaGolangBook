package codegen

type locVarInfo struct {
	prev    *locVarInfo //上一层
	name    string      //局部变量名字
	scopeLv int         //作用域层级
	slot    int         //绑定的索引
	capture bool        //局部变量是否被闭包不火
}

type funcInfo struct {
	constants map[interface{}]int    //常量表
	usedRegs  int                    //已经使用的寄存器
	maxRegs   int                    //最大寄存器
	scopeLv   int                    //作用域层级
	locVars   []*locVarInfo          //内部申明的全部局部变量
	locNames  map[string]*locVarInfo //当前生效的局部变量
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
func (self *funcInfo) enterScope() {
	self.scopeLv++
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
	self.scopeLv--
	for _, locVar := range self.locNames {
		if locVar.scopeLv > self.scopeLv { //离开作用域
			self.removeLocVar(locVar)
		}
	}
}
