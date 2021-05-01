package state

import . "LuaGo/api"

type luaStack struct {
	slots   []luaValue       //栈存放值
	top     int              //栈顶索引
	state   *luaState        //虚拟机
	closure *closure         //闭包
	varargs []luaValue       //实现变长参数
	openuvs map[int]*upvalue //捕获的局部变量(Open状态)
	pc      int              //内部指令
	prev    *luaStack        //形成链表节点用
}

//创建指定长度的lua stack
func newLuaStack(size int, state *luaState) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		top:   0,
		state: state,
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

//推入N个数值
func (self *luaStack) pushN(vals []luaValue, n int) {
	nVals := len(vals)
	if n < 0 {
		n = nVals
	}
	for i := 0; i < n; i++ {
		if i < nVals {
			self.push(vals[i])
		} else {
			self.push(nil)
		}

	}
}

//把索引转换成绝对索引
//一次性弹出N个值
func (self *luaStack) popN(n int) []luaValue {
	vals := make([]luaValue, n)
	for i := n - 1; i >= 0; i-- {
		vals[i] = self.pop()
	}
	return vals
}

func (self *luaStack) absIndex(idx int) int {
	//超出栈范围了
	if idx >= 0 || idx <= LUA_REGISTRYINDEX {
		return idx
	}

	return idx + self.top + 1
}

//索引是否有效  lua index起始位置1  最大 top
func (self *luaStack) isValid(idx int) bool {
	//upvalues
	if idx < LUA_REGISTRYINDEX {
		uvIdx := LUA_REGISTRYINDEX - idx - 1
		c := self.closure
		return c != nil && uvIdx < len(c.upvals)
	}

	//等于伪索引认定为有效索引
	if idx == LUA_REGISTRYINDEX {
		return true
	}
	absIdx := self.absIndex(idx)
	return absIdx > 0 && absIdx <= self.top
}

//从栈中取出某个index  go是0开始
func (self *luaStack) get(idx int) luaValue {
	//upvalues
	if idx < LUA_REGISTRYINDEX {
		uvIdx := LUA_REGISTRYINDEX - idx - 1
		c := self.closure
		if c == nil || uvIdx >= len(c.upvals) {
			return nil
		}
		return *(c.upvals[uvIdx].val)
	}

	//如果是伪装索引 则返回注册表
	if idx == LUA_REGISTRYINDEX {
		return self.state.registry
	}
	absIdx := self.absIndex(idx)
	if absIdx > 0 && absIdx <= self.top {
		return self.slots[absIdx-1]
	}
	return nil
}

//往栈里写入某个值  索引无效则panic
func (self *luaStack) set(idx int, val luaValue) {
	//upvalues
	if idx < LUA_REGISTRYINDEX {
		uvIdx := LUA_REGISTRYINDEX - idx - 1
		c := self.closure
		if c != nil && uvIdx < len(c.upvals) {
			*(c.upvals[uvIdx].val) = val
		}
		return
	}

	if idx == LUA_REGISTRYINDEX {
		self.state.registry = val.(*luaTable)
		return
	}

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




