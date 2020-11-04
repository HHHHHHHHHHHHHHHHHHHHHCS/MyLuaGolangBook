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
	self.setTable(t, k, v, false)
}

//存入表数据
//__newindex  t[k]=v t不是表或者k在表中不存在   可能会递归触发
//t 有可能是表 也有可能是函数
//raw true 代表不用元表 直接暴力存入
func (self *luaState) setTable(t, k, v luaValue, raw bool) {
	if tbl, ok := t.(*luaTable); ok {
		if raw || tbl.get(k) != nil || !tbl.hasMetafield("__newindex") {
			tbl.put(k, v)
			return
		}
	}
	if !raw {
		if mf := getMetafield(t, "__newindex", self); mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				self.setTable(x, k, v, false)
				return
			case *closure:
				self.stack.push(mf)
				self.stack.push(t)
				self.stack.push(k)
				self.stack.push(v)
				self.Call(3, 0)
				return
			}
		}
	}
	panic("index error")
}

//存入表key位置(string) 数据
func (self *luaState) SetField(idx int, k string) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self.setTable(t, k, v, false)
}

//存入表key位置(int) 数据
func (self *luaState) SetI(idx int, i int64) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self.setTable(t, i, v, false)
}
