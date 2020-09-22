package state

import (
	. "LuaGo/api"
	"LuaGo/binchunk"
)

//closure 闭包
type closure struct {
	proto  *binchunk.Prototype //lua closure
	goFunc GoFunction          //go closure
}

func newLuaClosure(proto *binchunk.Prototype) *closure {
	return &closure{proto: proto}
}

func newGoClosure(f GoFunction) *closure {
	return &closure{goFunc: f}
}
