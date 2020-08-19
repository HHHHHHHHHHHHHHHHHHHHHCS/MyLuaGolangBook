package state

import "math"
import "LuaGo/number"

type luaTable struct {
	arr  []luaValue
	_map map[luaValue]luaValue //字典  字段名不能跟struct 重复
}

func newLuaTable(nArr, nRec int) *luaTable {
	t := &luaTable{}
	if nArr > 0 {
		//第二个参数 切片长度 第三个参数 预留的空间
		t.arr = make([]luaValue, 0, nArr)
	}
	if nRec > 0 {
		t._map = make(map[luaValue]luaValue, nRec)
	}
	return t
}

func (self *luaTable) get(key luaValue) luaValue {
	key = _floatToInteger(key)
}

func _floatToInteger(key luaValue) luaValue {
	if f, ok := key.(float64); ok {
		if i, ok := number.FloatToInteger(f); ok {
			return i
		}
	}
	return key
}
