package state

type luaStack struct {
	slots   []luaValue //栈存放值
	top     int        //栈顶索引
	prev    *luaStack  //形成链表节点用
	closure *closure   //闭包
	varargs []luaValue //实现变长参数
	pc      int        //内部指令
}

//创建指定长度的lua stack
func newLuaStack(size int) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		top:   0,
	}
}

//判断容器是否可以容纳n  如果不行则扩容
func (self *luaStack) check(n int) {
	free := len(self.slots) - self.top
	for i := free; i < n; i++ {
		self.slots = append(self.slots, nil)
	}
}

//将值推入栈顶  如果溢出 则panic抛出
func (self *luaStack) push(val luaValue) {
	if self.top == len(self.slots) {
		panic("stack overflow!")
	}

	self.slots[self.top] = val
	self.top++
}

//从栈顶弹出一个值  如果栈是空的 则panic抛出
func (self *luaStack) pop() luaValue {
	if self.top < 1 {
		panic("stack underflow!")
	}

	self.top--
	val := self.slots[self.top]
	self.slots[self.top] = nil
	return val
}

//把索引转换成绝对索引
func (self *luaStack) absIndex(idx int) int {
	if idx >= 0 {
		return idx
	}
	return idx + self.top + 1
}

//索引是否有效  lua index起始位置1  最大 top
func (self *luaStack) isValid(idx int) bool {
	absIdx := self.absIndex(idx)
	return absIdx > 0 && absIdx <= self.top
}

//从栈中取出某个index  go是0开始
func (self *luaStack) get(idx int) luaValue {
	absIdx := self.absIndex(idx)
	if absIdx > 0 && absIdx <= self.top {
		return self.slots[absIdx-1]
	}
	return nil
}

//往栈里写入某个值  索引无效则panic
func (self *luaStack) set(idx int, val luaValue) {
	absIdx := self.absIndex(idx)
	if absIdx > 0 && absIdx <= self.top {
		self.slots[absIdx-1] = val
		return
	}
	panic("invalid index!")
}

//翻转
func (self *luaStack) reverse(from, to int) {
	slots := self.slots
	for from < to {
		slots[from], slots[to] = slots[to], slots[from]
		from++
		to--
	}
}
