package state

import . "LuaGo/api"

//比较操作 选择栈里的两个位置进行比较 不会修改栈的状态
func (self *luaState) Compare(idx1, idx2 int, op CompareOp) bool {
	a := self.stack.get(idx1)
	b := self.stack.get(idx2)

	switch op {
	case LUA_OPEQ:
		return _eq(a, b)
	case LUA_OPLT:
		return _lt(a, b)
	case LUA_OPLE:
		return _le(a, b)
	default:
		panic("invalid compare op!")
	}
}

/*
	判断相等
	只有两个相同类型才能比较  int64 和 float 要进行转换判断
	如果是 自定义数据 则根据地址判断
*/
func _eq(a, b luaValue) bool {
	switch x := a.(type) {
	case nil:
		return b == nil
	case bool:
		y, ok := b.(bool)
		return ok && x == y
	case string:
		y, ok := b.(string)
		return ok && x == y
	case int64:
		switch y := b.(type) {
		case int64:
			return x == y
		case float64:
			return float64(x) == y
		default:
			return false
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x == y
		case int64:
			return x == float64(y)
		default:
			return false
		}
	default:
		return a == b
	}
}

//a<b
func _lt(a, b luaValue) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x < y
		}
	case int64:
		{
			switch y := b.(type) {
			case int64:
				return x < y
			case float64:
				return float64(x) < y
			}
		}
	case float64:
		{
			switch y := b.(type) {
			case float64:
				return x < y
			case int64:
				return x < float64(y)
			}
		}
	}
	panic("comparison error!")
}

//a<=b
func _le(a, b luaValue) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x <= y
		}
	case int64:
		{
			switch y := b.(type) {
			case int64:
				return x <= y
			case float64:
				return float64(x) <= y
			}
		}
	case float64:
		{
			switch y := b.(type) {
			case float64:
				return x <= y
			case int64:
				return x <= float64(y)
			}
		}
	}
	panic("comparison error!")
}
