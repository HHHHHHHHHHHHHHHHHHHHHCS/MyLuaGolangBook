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

	//如果可以转换为int  即尝试用index序号
	key = _floatToInteger(key)
	if idx, ok := key.(int64); ok && idx >= 1 {
		//如果在数组范围的内
		arrLen := int64(len(self.arr))
		if idx <= arrLen {
			//go lua  数组起点
			self.arr[idx-1] = val
			if idx == arrLen && val == nil {
				//删除末尾的nil
				self._shrinkArray()
			}
			return
		}
		//超出 范围+1 直接是追加到数组
		if idx == arrLen+1 {
			//尝试删除  从字典删除   因为序号之前可以当作key 加入到字典
			delete(self._map, key)
			if val != nil {
				self.arr = append(self.arr, val)
				self._expandArray()
			}
			return
		}
	}

	//序号不是超出数组范围  序号不是整数 等 存入字典
	if val != nil {
		if self._map == nil {
			self._map = make(map[luaValue]luaValue, 8)
		}
		self._map[key] = val
	} else {
		delete(self._map, key)
	}
}

//可能存在切片错误的问题
func (self *luaTable) _shrinkArray() {
	for i := len(self.arr) - 1; i >= 0; i-- {
		if self.arr[i] == nil {
			self.arr = self.arr[0:i]
		}
	}
}

func (self *luaTable) _expandArray() {
	for idx := int64(len(self.arr)) + 1; true; idx++ {
		if val, found := self._map[idx]; found {
			delete(self._map, idx)
			self.arr = append(self.arr, val)
		} else {
			break
		}
	}
}

func (self *luaTable) len() int {
	return len(self.arr)
}

