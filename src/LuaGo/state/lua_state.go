package state

type luaState struct {
	stack *luaStack
}

func New() *luaState {
	return &luaState{
		stack: newLuaStack(20),
	}
}

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
	for i := 0; i < n; i++ {
		self.stack.pop()
	}
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
func (self *luaState) Insert(idx int){
	self.Rotate(idx, 1)
}

//移除某个值 其余全部前挪一位
func (self *luaState) Remove(idx int) {
	self.Rotate(idx, -1)
	self.Pop()
}