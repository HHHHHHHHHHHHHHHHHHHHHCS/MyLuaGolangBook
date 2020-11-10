package state

import (
	. "LuaGo/api"
	"fmt"
)

//不使用原表 直接len
func (self *luaState) RawLen(idx int) uint {
	val := self.stack.get(idx)
	switch x := val.(type) {
	case string:
		return uint(len(x))
	case *luaTable:
		return uint(x.len())
	default:
		return 0
	}
}

//从栈顶获取信息
//把LUA的类型 转换成字符串
func (self *luaState) TypeName(tp LuaType) string {
	switch tp {
	case LUA_TNONE:
		return "no value"
	case LUA_TNIL:
		return "nil"
	case LUA_TBOOLEAN:
		return "boolean"
	case LUA_TNUMBER:
		return "number"
	case LUA_TSTRING:
		return "string"
	case LUA_TTABLE:
		return "table"
	case LUA_TFUNCTION:
		return "function"
	case LUA_TTHREAD:
		return "thread"
	default:
		return "userdata"
	}
}

//根据索引的值 返回类型
func (self *luaState) Type(idx int) LuaType {
	if self.stack.isValid(idx) {
		val := self.stack.get(idx)
		return typeOf(val)
	}
	return LUA_TNONE
}

func (self *luaState) IsNone(idx int) bool {
	return self.Type(idx) == LUA_TNONE
}

func (self *luaState) IsNil(idx int) bool {
	return self.Type(idx) == LUA_TNIL
}

func (self *luaState) IsNoneOrNil(idx int) bool {
	return self.Type(idx) <= LUA_TNIL
}

func (self *luaState) IsBoolean(idx int) bool {
	return self.Type(idx) == LUA_TBOOLEAN
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_istable
func (self *luaState) IsTable(idx int) bool {
	return self.Type(idx) == LUA_TTABLE
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isfunction
func (self *luaState) IsFunction(idx int) bool {
	return self.Type(idx) == LUA_TFUNCTION
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isthread
func (self *luaState) IsThread(idx int) bool {
	return self.Type(idx) == LUA_TTHREAD
}

//索引处的值是否是字符串  数字可以解析为字符串
func (self *luaState) IsString(idx int) bool {
	t := self.Type(idx)
	return t == LUA_TSTRING || t == LUA_TNUMBER
}

//引处的值是否是数字 或者 可以转换为数字
func (self *luaState) IsNumber(idx int) bool {
	_, ok := self.ToNumberX(idx)
	return ok
}

func (self *luaState) IsInteger(idx int) bool {
	val := self.stack.get(idx)
	_, ok := val.(int64)
	return ok
}

//取出一个指 如果不是布尔值 则进行转换
func (self *luaState) ToBoolean(idx int) bool {
	val := self.stack.get(idx)
	return convertToBoolean(val)
}

//进行普通的数字转换 如果不是数字 则返回0
func (self *luaState) ToNumber(idx int) float64 {
	n, _ := self.ToNumberX(idx)
	return n
}

//进行带转换结果的数字转换 如果不是数字 则返回 0 和 结果
func (self *luaState) ToNumberX(idx int) (float64, bool) {
	val := self.stack.get(idx)
	return convertToFloat(val)
}

//转换为整数 失败为 0
func (self *luaState) ToInteger(idx int) int64 {
	i, _ := self.ToIntegerX(idx)
	return i
}

//转换为整数 和  转换结果
func (self *luaState) ToIntegerX(idx int) (int64, bool) {
	val := self.stack.get(idx)
	return convertToInteger(val)
}

//转换为字符串 数字也可以转换为字符串 失败为空字符串
func (self *luaState) ToString(idx int) string {
	s, _ := self.ToStringX(idx)
	return s
}

//转换为字符串  和  转换结果
func (self *luaState) ToStringX(idx int) (string, bool) {
	val := self.stack.get(idx)
	switch x := val.(type) {
	case string:
		return x, true
	case int64, float64:
		s := fmt.Sprintf("%v", x)
		self.stack.set(idx, s) //注意这里会修改栈
		return s, true
	default:
		return "", false
	}
}

//检测是否是go闭包 拿到变量 是否是闭包 是否不为nil
func (self *luaState) IsGoFunction(idx int) bool {
	val := self.stack.get(idx)
	if c, ok := val.(*closure); ok {
		return c.goFunc != nil
	}
	return false
}

//从栈中返回goFunc
func (self *luaState) ToGoFunction(idx int) GoFunction {
	val := self.stack.get(idx)
	if c, ok := val.(*closure); ok {
		return c.goFunc
	}
	return nil
}
