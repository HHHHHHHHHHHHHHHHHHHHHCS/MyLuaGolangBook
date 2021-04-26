package state

import (
	. "LuaGo/api"
	"io/ioutil"
)

func (self *luaState) TypeName2(idx int) string {
	return self.TypeName(self.Type(idx))
}

func (self *luaState) CheckStack2(sz int, msg string) {
	if !self.CheckStack(sz) {
		if msg != "" {
			self.Error2("stack overflow (%s)", msg)
		} else {
			self.Error2("stack overflow")
		}
	}
}

func (self *luaState) Error2(fmt string, a ...interface{}) int {
	self.PushFString(fmt, a...)
	return self.Error()
}

func (self *luaState) LoadString(s string) int {
	return self.Load([]byte(s), s, "bt")
}

func (self *luaState) LoadFileX(filename, mode string) int {
	if data, err := ioutil.ReadFile(filename); err == nil {
		return self.Load(data, "@"+filename, mode)
	}
	return LUA_ERRFILE
}

func (self *luaState) LoadFile(filename string) int {
	return self.LoadFileX(filename, "bt")
}

func (self *luaState) DoString(str string) bool {
	return self.LoadString(str) == LUA_OK &&
		self.PCall(0, LUA_MULTRET, 0) == LUA_OK
}

func (self *luaState) DoFile(filename string) bool {
	return self.LoadFile(filename) == LUA_OK &&
		self.PCall(0, LUA_MULTRET, 0) == LUA_OK
}

func (self *luaState) ArgError(arg int, extraMsg string) int {
	return self.Error2("bad argument $%d (%s)", arg, extraMsg)
}

func (self *luaState) ArgCheck(cond bool, arg int, extraMsg string) {
	if !cond {
		self.ArgError(arg, extraMsg)
	}
}

