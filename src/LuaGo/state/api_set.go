package state

import . "LuaGo/api"

//把栈顶的值写到全局注册表里面 name是key
func (self *luaState) SetGlobal(name string) {
	t := self.registry.get(LUA_RIDX_GLOBALS)
	v := self.stack.pop()
	self.setTable(t, name, v)
}

//把栈顶的go函数 写到表里面
func (self *luaState) Register(name string, f GoFunction) {
	self.PushGoFunction(f)
	self.SetGlobal(name)
}

//查找idx 的表  栈顶弹出两个值当作 k,v  存入表
func (self *luaState) SetTable(idx int) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	k := self.stack.pop()
	self.setTable(t, k, v)
}

//存入表数据
func (self *luaState) setTable(t, k, v luaValue) {
	if tbl, ok := t.(*luaTable); ok {
		tbl.put(k, v)
		return
	}
	panic("not a table")
}

//存入表key位置(string) 数据
func (self *luaState) SetField(idx int, k string) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self.setTable(t, k, v)
}

//存入表key位置(int) 数据
func (self *luaState) SetI(idx int, i int64) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self.setTable(t, i, v)
}
