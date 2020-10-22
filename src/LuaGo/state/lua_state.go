package state

import . "LuaGo/api"

type luaState struct {
	registry *luaTable //注册表
	stack    *luaStack
}

//Lua栈初始容量
func New() *luaState {
	registry := newLuaTable(0, 0)
	//每个lua解释器都有自己的注册表
	registry.put(LUA_RIDX_GLOBALS, newLuaTable(0, 0)) //全局变量

	ls := &luaState{registry: registry}
	ls.pushLuaStack(newLuaStack(LUA_MINSTACK, ls))
	return ls
}

func (self *luaState) pushLuaStack(stack *luaStack) {
	stack.prev = self.stack
	self.stack = stack
}

func (self *luaState) popLuaStack() {
	stack := self.stack
	self.stack = stack.prev
	stack.prev = nil
}
