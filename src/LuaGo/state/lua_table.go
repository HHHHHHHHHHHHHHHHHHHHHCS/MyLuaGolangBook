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

//从map get value    如果key可以转换为int 则尝试
//go起始位置是0   lua起始位置是1
func (self *luaTable) get(key luaValue) luaValue {
	key = _floatToInteger(key)
	if idx, ok := key.(int64); ok {
		if idx >= 1 && idx <= int64(len(self.arr)) {
			return self.arr[idx-1]
		}
	}
	return self._map[key]
}

//key float 尝试转换 int
func _floatToInteger(key luaValue) luaValue {
	if f, ok := key.(float64); ok {
		if i, ok := number.FloatToInteger(f); ok {
			return i
		}
	}
	return key
}

func (self *luaTable) put(key, val luaValue) {
	//不能是nil
	if key == nil {
		panic("table index is nil!")
	}
	// 如果是float 就不能是nan
	if f, ok := key.(float64); ok && math.IsNaN(f) {
		panic("table index is Nan!")
	}
	key = _floatToInteger(key)
	if idx, ok := key.(int64); ok && idx >= 1 {
		arrLen := int64(len(self.arr))
		if idx <= arrLen {
			if idx == arrLen && val == nil {
				self._ //TODO:
			}
			return
		}
		//超出范围
		if idx == arrLen+1 {
			//尝试删除    因为key 之前可能非序号的形式存入  现在是序号模式  所以要先删除
			delete(self._map, key)
			if val != nil {
				self.arr = append(self.arr, val)
				self._exp //TODO:
			}
			return
		}
	}
	if val != nil {
		if self._map == nil {
			self._map = make(map[luaValue]luaValue, 8)
		} else {
			delete(self._map, key)
		}
	}
}
