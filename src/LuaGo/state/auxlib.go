package state

import (
	. "LuaGo/api"
	"LuaGo/stdlib"
	"fmt"
	"io/ioutil"
)

func (self *luaState) Error2(fmt string, a ...interface{}) int {
	self.PushFString(fmt, a...)
	return self.Error()
}

func (self *luaState) ArgError(arg int, extraMsg string) int {
	return self.Error2("bad argument #%d (%s)", arg, extraMsg)
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

func (self *luaState) CheckInteger(arg int) int64 {
	i, ok := self.ToIntegerX(arg)
	if !ok {
		self.intError(arg)
	}
	return i
}

func (self *luaState) CheckNumber(arg int) float64 {
	f, ok := self.ToNumberX(arg)
	if !ok {
		self.tagError(arg, LUA_TNUMBER)
	}
	return f
}

func (self *luaState) CheckString(arg int) string {
	s, ok := self.ToStringX(arg)
	if !ok {
		self.tagError(arg, LUA_TSTRING)
	}
	return s
}

func (self *luaState) OptInteger(arg int, def int64) int64 {
	if self.IsNoneOrNil(arg) {
		return def
	}
	return self.CheckInteger(arg)
}

//检查参数是否是选定的类型 如果是则返回值  否则返回默认值
func (self *luaState) OptNumber(arg int, def float64) float64 {
	if self.IsNoneOrNil(arg) {
		return def
	}
	return self.CheckNumber(arg)
}

func (self *luaState) OptString(arg int, def string) string {
	if self.IsNoneOrNil(arg) {
		return def
	}
	return self.CheckString(arg)
}

func (self *luaState) DoFile(filename string) bool {
	return self.LoadFile(filename) != LUA_OK ||
		self.PCall(0, LUA_MULTRET, 0) != LUA_OK
}

func (self *luaState) DoString(str string) bool {
	return self.LoadString(str) != LUA_OK ||
		self.PCall(0, LUA_MULTRET, 0) != LUA_OK
}

func (self *luaState) LoadFile(filename string) int {
	return self.LoadFileX(filename, "bt")
}

func (self *luaState) LoadFileX(filename, mode string) int {
	if data, err := ioutil.ReadFile(filename); err == nil {
		return self.Load(data, "@"+filename, mode)
	}
	return LUA_ERRFILE
}

func (self *luaState) LoadString(s string) int {
	return self.Load([]byte(s), s, "bt")
}

func (self *luaState) TypeName2(idx int) string {
	return self.TypeName(self.Type(idx))
}

func (self *luaState) Len2(idx int) int64 {
	self.Len(idx)
	i, isNum := self.ToIntegerX(-1)
	if !isNum {
		self.Error2("object length is not an integer")
	}
	self.Pop(1)
	return i
}

func (self *luaState) ToString2(idx int) string {
	if self.CallMeta(idx, "__tostring") {
		if !self.IsString(-1) {
			self.Error2("'__tostring' must return a string")
		}
	} else {
		switch self.Type(idx) {
		case LUA_TNUMBER:
			if self.IsInteger(idx) {
				self.PushString(fmt.Sprintf("%d", self.ToInteger(idx)))
			} else {
				self.PushString(fmt.Sprintf("%g", self.ToNumber(idx)))
			}
			break
		case LUA_TSTRING:
			self.PushValue(idx)
			break
		case LUA_TBOOLEAN:
			if self.ToBoolean(idx) {
				self.PushString("true")
			} else {
				self.PushString("false")
			}
			break
		case LUA_TNIL:
			self.PushString("nil")
			break
		default:
			tt := self.GetMetafield(idx, "__name")
			var kind string
			if tt == LUA_TSTRING {
				kind = self.CheckString(-1)
			} else {
				kind = self.TypeName2(idx)
			}
			self.PushString(fmt.Sprintf("%s: %p", kind, self.ToPointer(idx)))
			if tt != LUA_TNIL {
				self.Remove(-2)
			}
		}
	}
	return self.CheckString(-1)
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

//查看 元表是否有这个属性
func (self *luaState) GetMetafield(obj int, event string) LuaType {
	if !self.GetMetatable(obj) {
		return LUA_TNIL
	}

	self.PushString(event)
	tt := self.RawGet(-2)
	if tt == LUA_TNIL {
		self.Pop(2)
	} else {
		self.Remove(-2)
	}
	return tt
}

func (self *luaState) CallMeta(obj int, event string) bool {
	obj = self.AbsIndex(obj)
	if self.GetMetafield(obj, event) == LUA_TNIL {
		return false
	}

	self.PushValue(obj)
	self.Call(1, 1)
	return true
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

func (self *luaState) NewLib(l FuncReg) {
	self.NewLibTable(l)
	self.SetFuncs(l, 0)
}

func (self *luaState) NewLibTable(l FuncReg) {
	self.CreateTable(0, len(l))
}

func (self *luaState) SetFuncs(l FuncReg, nup int) {
	self.CheckStack2(nup, "too many upvalues")
	for name, fun := range l {
		for i := 0; i < nup; i++ {
			self.PushValue(-nup)
		}

		self.PushGoClosure(fun, nup)
		self.SetField(-(nup + 2), name)
	}
	self.Pop(nup)
}

func (self *luaState) intError(arg int) {
	if self.IsNumber(arg) {
		self.ArgError(arg, "number has no integer representation")
	} else {
		self.tagError(arg, LUA_TNUMBER)
	}
}

func (self *luaState) tagError(arg int, tag LuaType) {
	self.typeError(arg, self.TypeName(LuaType(tag)))
}

func (self *luaState) typeError(arg int, tname string) int {
	var typeArg string
	if self.GetMetafield(arg, "__name") == LUA_TSTRING {
		typeArg = self.ToString(-1)
	} else if self.Type(arg) == LUA_TLIGHTUSERDATA {
		typeArg = "light userdata"
	} else {
		typeArg = self.TypeName2(arg)
	}
	msg := tname + " expected, got " + typeArg
	self.PushString(msg)
	return self.ArgError(arg, msg)
}