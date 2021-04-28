package stdlib

import . "LuaGo/api"

var baseFuncs = map[string]GoFunction{
	"print":        basePrint,
	"assert":       baseAssert,
	"error":        baseError,
	"select":       baseSelect,
	"ipairs":       baseIPairs,
	"next":         baseNext,
	"load":         baseLoad,
	"loadfile":     baseLoadFile,
	"dofile":       baseDoFile,
	"pcall":        basePCall,
	"xpcall":       baseXPCall,
	"getmetatable": baseGetMetatable,
	"setmetatable": baseSetMetatable,
	"rawequal":     baseRawEqual,
	"rawlen":       baseRawLen,
	"rawget":       baseRawGet,
	"rawset":       baseRawSet,
	"type":         baseType,
	"tostring":     baseToString,
	"tonumber":     baseToNumber,
}

func OpenBaseLib(ls LuaState) int {
	//open lib into global table
	ls.PushGlobalTable()
	ls.SetFuncs(baseFuncs, 0)
	//set global _G
	ls.PushValue(-1)
	ls.SetField(-2, "_G")
	//set global _Version
	ls.PushString("Lua 5.3")
	ls.SetField(-2, "_VERSION")
	return 1
}

func basePrint(ls LuaState) {
	//todo:
}

func baseSelect(ls LuaState) int {
	n := int64(ls.GetTop())
	if ls.Type(1) == LUA_TSTRING && ls.CheckString(1) == "#" {
		//#返回参数个数
		ls.PushInteger(n - 1)
		return 1
	} else {
		i := ls.CheckInteger(1)
		if i < 0 {
			i = n + i
		} else if i > n {
			i = n
		}
		ls.ArgCheck(1 <= i, 1, "index out of range")
		//把适当数量的参数作为返回值 就好了
		return int(n - 1)
	}

}
