package state

//获取当前程序计数器
func (self *luaState) PC() int {
	return self.pc
}

//当前的程序计数器+n
func (self *luaState) AddPC(n int) {
	self.pc += n
}

//从指令表中读取出指令  指令计数器+1  下次获取的时候就是下一条指令
func (self *luaState) Fetch() uint32 {
	i := self.proto.Code[self.pc]
	self.pc++
	return i
}

//从常量表中读取一个常量 压入当前的栈顶
func (self *luaState) GetConst(idx int) {
	c := self.proto.Constants[idx]
	self.stack.push(c)
}

//如果RK 大于255 则视为 取出常量值  否则 把 索引压入栈顶
//常用于 如iABC的OpArgK 高位加1 取值 高位是0 压入索引
func (self *luaState) GetRK(rk int) {
	if rk > 0xFF { //constant
		self.GetConst(rk & 0xFF)
	} else {
		//lua的索引是从零开始的
		//但是luaAPI 的栈是从1开始的 使用的时候要加一
		self.PushValue(rk + 1)
	}

}