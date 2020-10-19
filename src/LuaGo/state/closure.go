package state

import (
	. "LuaGo/api"
	. "LuaGo/binchunk"
)

type upvalue struct {
	val *luaValue
}

//closure 闭包
type closure struct {
	proto  *Prototype //lua closure
	goFunc GoFunction //go closure
	upvals []*upvalue //values
}

func newLuaClosure(proto *Prototype) *closure {
	c := &closure{proto: proto}
	if nUpvals := len(proto.Upvalues); nUpvals > 0 {
		c.upvals = make([]*upvalue, nUpvals)
	}
	return c
}

//go的需要指明闭包的upvalues的数量
func newGoClosure(f GoFunction, nUpvals int) *closure {
	c := &closure{goFunc: f}
	if nUpvals > 0 {
		c.upvals = make([]*upvalue, nUpvals)
	}
	return c
}
