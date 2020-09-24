package state

import (
	. "LuaGo/api"
)

func (self *luaState) CreateTable(nArr, nRec int) {
	t := newLuaTable(nArr, nRec)
	self.stack.push(t)
}

func (self *luaState) NewTable() {
	self.CreateTable(0, 0)
}

//从idx找出表  索引在栈顶
func (self *luaState) GetTable(idx int) LuaType {
	t := self.stack.get(idx)
	k := self.stack.pop()
	return self.getTable(t, k)
}

//查找表  查找key得到value  通过 idx 得到表    k 是key 可能会自动转换为int 去查找
func (self *luaState) GetField(idx int, k string) LuaType {
	t := self.stack.get(idx)
	return self.getTable(t, k)
}

//查找表  但是这时候key 是int 效率高点 不用判断和转换
func (self *luaState) GetI(idx int, i int64) LuaType {
	t := self.stack.get(idx)
	return self.getTable(t, i)
}

//t 是 table  get(index)的value  把值放入栈顶  返回val.(typeOf)
func (self *luaState) getTable(t, k luaValue) LuaType {
	if tbl, ok := t.(*luaTable); ok {
		v := tbl.get(k)
		self.stack.push(v)
		return typeOf(v)
	}
	panic("not a table!")
}