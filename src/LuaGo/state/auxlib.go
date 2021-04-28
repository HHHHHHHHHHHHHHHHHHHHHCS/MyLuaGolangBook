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

func (self *luaState) CheckAny(arg int) {
	if self.Type(arg) == LUA_TNONE {
		self.ArgError(arg, "value expected")
	}
}

func (self *luaState) CheckType(arg int, t LuaType) {
	if self.Type(arg) != t {
		self.tagError(arg, t)
	}
}

func (self *luaState) CheckNumber(arg int) float64 {
	f, ok := self.ToNumberX(arg)
	if !ok {
		self.tagError(arg, LUA_TNUMBER)
	}
	return f
}

//检查参数是否是选定的类型 如果是则返回值  否则返回默认值
func (self *luaState) OptNumber(arg int, def float64) float64 {
	if self.IsNoneOrNil(arg) {
		return def
	}
	return self.CheckNumber(arg)
}

func (self *luaState) OpenLibs() {
	libs := map[string]GoFunction{
		"_G": stdlib.OpenBaseLib,
		//TODO:
	}

	for name, fun := range libs {
		self.RequireF(name, fun, true)
		self.Pop(1)
	}
}

func (self *luaState) RequireF(modname string, openf GoFunction, glb bool) {
	self.GetSubTable(LUA_REGISTRYINDEX, "_LOADED")
	//loaded modname
	self.GetField(-1, modname)

	//package not already loaded?
	if !self.ToBoolean(-1) {
		self.Pop(1) //remove field
		self.PushGoFunction(openf)
		self.PushString(modname)   //arg to open func
		self.Call(1, 1)            //call openf to open module
		self.PushValue(-1)         //make copy of module(call result)
		self.SetField(-3, modname) //_LOADED[modname] = module
	}

	self.Remove(-2) //remove _LOADED table
	if glb {
		self.PushValue(-1)      //copy of module
		self.SetGlobal(modname) //_G[modname] = module
	}
}

func (self *luaState) GetSubTable(idx int, fname string) bool {
	if self.GetField(idx, fname) == LUA_TTABLE {
		return true //table already there
	}

	self.Pop(1) //remove previous result
	idx = self.stack.absIndex(idx)
	self.NewTable()
	self.PushValue(-1)        //copy to be left at top
	self.SetField(idx, fname) //assign new table to field
	return false              //false ,because did not find table there
}


