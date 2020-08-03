package state

import "LuaGo/binchunk"

type luaState struct {
	stack *luaStack
	proto *binchunk.Prototype
	pc    int
}

//Lua栈初始容量 和 函数原型以初始化proto字段
func New(stackSize int, proto *binchunk.Prototype) *luaState {
	return &luaState{
		stack: newLuaStack(stackSize),
		proto: proto,
		pc:    0,
	}
}
