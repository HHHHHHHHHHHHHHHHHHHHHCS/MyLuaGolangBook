package state

//获取当前程序计数器
func (self *luaState) PC() int {
	return self.stack.pc
}

//当前的程序计数器+n
func (self *luaState) AddPC(n int) {
	self.stack.pc += n
}

//从指令表中读取出指令  指令计数器+1  下次获取的时候就是下一条指令
func (self *luaState) Fetch() uint32 {
	i := self.stack.closure.proto.Code[self.stack.pc]
	self.stack.pc++
	return i
}

//从常量表中读取一个常量 压入当前的栈顶
func (self *luaState) GetConst(idx int) {
	c := self.stack.closure.proto.Constants[idx]
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

//获取寄存器数量
func (self *luaState) RegisterCount() int {
	return int(self.stack.closure.proto.MaxStackSize)
}

//读取可变参数 放入栈顶
func (self *luaState) LoadVararg(n int) {
	if n < 0 {
		n = len(self.stack.varargs)
	}

	self.stack.check(n)
	self.stack.pushN(self.stack.varargs, n)
}

//获取函数原型(子函数) 生成闭包
func (self *luaState) LoadProto(idx int) {
	stack := self.stack
	subProto := stack.closure.proto.Protos[idx]
	closure := newLuaClosure(subProto)
	self.stack.push(closure)
	for i, uvInfo := range subProto.Upvalues {
		uvIdx := int(uvInfo.Idx)
		//Instack == 1  如果是局部变量
		if uvInfo.Instack == 1 {
			//只用访问当前函数的局部变量
			if stack.openuvs == nil {
				stack.openuvs = map[int]*upvalue{}
			}
			//如果还在缓存上(OPEN状态) 直接获取
			if openuv, found := stack.openuvs[uvIdx]; found {
				closure.upvals[i] = openuv
			} else {
				//在其他地方(CLOSE状态) 从栈上获取
				closure.upvals[i] = &upvalue{&stack.slots[uvIdx]}
				stack.openuvs[uvIdx] = closure.upvals[i]
			}
		} else {
			//如果是外部函数 则传递
			closure.upvals[i] = stack.closure.upvals[uvIdx]
		}
	}
}

func (self *luaState) CloseUpvalues(a int) {
	for i, openuv := range self.stack.openuvs {
		if i >= a-1 {
			//因为可能upvalue可能还存在引用的情况
			//先从寄存器拷贝出来Lua值 更新upvalue 为拷贝值
			val := *openuv.val
			openuv.val = &val
			//再删除 形成闭合状态
			delete(self.stack.openuvs, i)
		}
	}
}
