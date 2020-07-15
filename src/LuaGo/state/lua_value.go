package state

import (
	. "LuaGo/api"
)

type luaValue interface{}

func typeof(val luaValue) LuaType {
	//x.(type) 只能在switch 中使用
	switch val.(type) {
	case nil:
		return LUA_TNIL
	case bool:
		return LUA_TBOOLEAN
	case int64:
		return LUA_TNUMBER
	case float64:
		return LUA_TNUMBER
	case string:
		return LUA_TSTRING
	default:
		panic("todo!")
	}
}

//把类型转换为 boolean  nil->false  bool->bool  default->true
func convertToBoolean(val luaValue) bool {
	switch x := val.(type) {
	case nil:
		return false
	case bool:
		return x
	default:
		return true
	}
}
