package state

import (
	. "LuaGo/api"
	"fmt"
)

//lua值由外部推入栈顶
func (self *luaState) PushNil() {
	self.stack.push(nil)
}

func (self *luaState) PushBoolean(b bool) {
	self.stack.push(b)
}

func (self *luaState) PushInteger(n int64) {
	self.stack.push(n)
}
func (self *luaState) PushNumber(n float64) {
	self.stack.push(n)
}

func (self *luaState) PushString(s string) {
	self.stack.push(s)
}

func (self *luaState) PushFString(fmtStr string, a ...interface{}) {
	str := fmt.Sprintf(fmtStr, a...)
	self.stack.push(str)
}

func (self *luaState) PushGoFunction(f GoFunction) {
	self.stack.push(newGoClosure(f, 0))
}

func (self *luaState) PushGoClosure(f GoFunction, n int) {
	closure := newGoClosure(f, n)
	for i := n; i > 0; i-- {
		val := self.stack.pop()
		//索引从0开始
		closure.upvals[i-1] = &upvalue{&val}
	}
	self.stack.push(closure)
}

func (self *luaState) PushGlobalTable() {
	global := self.registry.get(LUA_RIDX_GLOBALS)
	self.stack.push(global)
}

func(self *luaState) PushThread() bool{
	self.stack.push(self)
	return self.isMainThread()
}