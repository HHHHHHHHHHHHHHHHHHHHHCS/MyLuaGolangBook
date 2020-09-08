package state

import "LuaGo/binchunk"

//closure 闭包
type closure struct {
	proto *binchunk.Prototype
}

func newLuaClosure(proto *binchunk.Prototype) *closure {
	return &closure{proto: proto}
}


