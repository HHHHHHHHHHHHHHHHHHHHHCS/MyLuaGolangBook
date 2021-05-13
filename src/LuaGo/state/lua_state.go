package state

import . "LuaGo/api"

type luaState struct {
	registry *luaTable //注册表
	stack    *luaStack
	coStatus int
	coCaller *luaState
	coChan   chan int
}

//Lua栈初始容量
func New() *luaState {
	ls := &luaState{}

	//每个lua解释器都有自己的注册表
	registry := newLuaTable(8, 0)
	registry.put(LUA_RIDX_MAINTHREAD, ls)              //注册主线程索引
	registry.put(LUA_RIDX_GLOBALS, newLuaTable(0, 20)) //全局变量

	ls.registry = registry
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

func (self *luaState) isMainThread() bool {
	return self.registry.get(LUA_RIDX_MAINTHREAD) == self
}
