package state

import . "LuaGo/api"

//lua值由外部推入栈顶
func (self *luaState) PushNil()                    { self.stack.push(nil) }
func (self *luaState) PushBoolean(b bool)          { self.stack.push(b) }
func (self *luaState) PushInteger(n int64)         { self.stack.push(n) }
func (self *luaState) PushNumber(n float64)        { self.stack.push(n) }
func (self *luaState) PushString(s string)         { self.stack.push(s) }
func (self *luaState) PushGoFunction(f GoFunction) { self.stack.push(newGoClosure(f, 0)) }
func (self *luaState) PushGlobalTable() {
	global := self.registry.get(LUA_RIDX_GLOBALS)
	self.stack.push(global)
}

func (self *luaState) GetGlobal(name string) LuaType {
	t := self.registry.get(LUA_RIDX_GLOBALS)
	return self.getTable(t, name, false)
}

func (self *luaState) PushGoClosure(f GoFunction, n int) {
	closure := newGoClosure(f, n)
	for i := n; i > 0; i-- {
		val := self.stack.pop()
		//索引从0开始
		closure.upvals[n-1] = &upvalue{&val}
	}
	self.stack.push(closure)
}
