package state

//返回栈顶索引
func (self *luaState) GetTop() int {
	return self.stack.top
}

//返回绝对索引
func (self *luaState) AbsIndex(idx int) int {
	return self.stack.absIndex(idx)
}

//检查扩容  如果过小则扩容
func (self *luaState) CheckStack(n int) bool {
	self.stack.check(n)
	return true //never fails
}

//栈顶弹出N个值
func (self *luaState) Pop(n int) {
	self.SetTop(-n - 1)
}

//把拷贝某个值到某个地方
func (self *luaState) Copy(fromIdx, toIdx int) {
	val := self.stack.get(fromIdx)
	self.stack.set(toIdx, val)
}

//把指定索引的值压入栈顶
func (self *luaState) PushValue(idx int) {
	val := self.stack.get(idx)
	self.stack.push(val)
}

//溢出栈顶 替换一个新的值
func (self *luaState) Replace(idx int) {
	val := self.stack.pop()
	self.stack.set(idx, val)
}

//把栈顶的值插入到某个地方
func (self *luaState) Insert(idx int) {
	self.Rotate(idx, 1)
}

//移除某个值 其余全部前挪一位
func (self *luaState) Remove(idx int) {
	self.Rotate(idx, -1)
	self.Pop(1)
}

//将[idx,top]内的值朝栈顶旋转n个方向 如果n是负数 则旋转到栈底
//lua进行是那次旋转
func (self *luaState) Rotate(idx, n int) {
	t := self.stack.top - 1
	p := self.stack.absIndex(idx) - 1
	var m int
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}
	self.stack.reverse(p, m)
	self.stack.reverse(m+1, t)
	self.stack.reverse(p, t)

}

//将栈顶设为指定值  如果小于当前索引 则相当于弹出
func (self *luaState) SetTop(idx int) {
	newTop := self.stack.absIndex(idx)
	if newTop < 0 {
		panic("stack underflow!")
	}

	n := self.stack.top - newTop
	if n > 0 {
		for i := 0; i < n; i++ {
			self.stack.pop()
		}
	} else if n < 0 {
		for i := 0; i > n; i-- {
			self.stack.push(nil)
		}
	}
}
