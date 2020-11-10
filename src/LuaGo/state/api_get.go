package state

import (
	. "LuaGo/api"
)

func (self *luaState) NewTable() {
	self.CreateTable(0, 0)
}

func (self *luaState) CreateTable(nArr, nRec int) {
	t := newLuaTable(nArr, nRec)
	self.stack.push(t)
}

//从idx找出表  索引在栈顶
func (self *luaState) GetTable(idx int) LuaType {
	t := self.stack.get(idx)
	k := self.stack.pop()
	return self.getTable(t, k, false)
}

//查找表  查找key得到value  通过 idx 得到表    k 是key 可能会自动转换为int 去查找
func (self *luaState) GetField(idx int, k string) LuaType {
	t := self.stack.get(idx)
	return self.getTable(t, k, false)
}

//查找表  但是这时候key 是int 效率高点 不用判断和转换
func (self *luaState) GetI(idx int, i int64) LuaType {
	t := self.stack.get(idx)
	return self.getTable(t, i, false)
}

//不使用元表  直接Get
func (self *luaState) RawGet(idx int) LuaType {
	t := self.stack.get(idx)
	k := self.stack.pop()
	return self.getTable(t, k, true)
}

//不使用元表  直接GetI
func (self *luaState) RawGetI(idx int, i int64) LuaType {
	t := self.GetTable(idx)
	return self.getTable(t, i, true)
}

func (self *luaState) GetGlobal(name string) LuaType {
	t := self.registry.get(LUA_RIDX_GLOBALS)
	return self.getTable(t, name, false)
}

func (self *luaState) GetMetatable(idx int) bool {
	val := self.stack.get(idx)

	if mt := getMetatable(val, self); mt != nil {
		self.stack.push(mt)
		return true
	} else {
		return false
	}
}

//t 是 table  get(index)的value  把值放入栈顶  返回val.(typeOf)
//__index t[k] 如果t不是表或者k不再表中不存在 则触发
//raw true代表 暴力直接table.get 忽略元方法
func (self *luaState) getTable(t, k luaValue, raw bool) LuaType {
	if tbl, ok := t.(*luaTable); ok {
		v := tbl.get(k)
		if raw || v != nil || !tbl.hasMetafield("__index") {
			self.stack.push(v)
			return typeOf(v)
		}
	}
	if !raw {
		//如果是table 则继续递归查找
		//如果是closure方法  则执行方法获取返回值
		if mf := getMetafield(t, "__index", self); mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				return self.getTable(x, k, false)
			case *closure:
				self.stack.push(mf)
				self.stack.push(t)
				self.stack.push(k)
				self.Call(2, 1)
				v := self.stack.get(-1)
				return typeOf(v)
			}
		}
	}
	panic("index error!")
}
