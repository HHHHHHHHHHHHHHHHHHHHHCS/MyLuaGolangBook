package state

import (
	"LuaGo/number"
	. "LuaGo/api"
	"math"
)

//Arith 执行算数和按位运算
//对于二元运算 弹出栈顶两个值 再压入运算结果
//对于一元运算 弹出栈顶一个值 在把结果压入

var (
	iadd  = func(a, b int64) int64 { return a + b }
	fadd  = func(a, b float64) float64 { return a + b }
	isub  = func(a, b int64) int64 { return a - b }
	fsub  = func(a, b float64) float64 { return a - b }
	imul  = func(a, b int64) int64 { return a * b }
	fmul  = func(a, b float64) float64 { return a * b }
	imod  = number.IMod //整数取模
	fmod  = number.FMod //浮点取模
	pow   = math.Pow
	div   = func(a, b float64) float64 { return a / b }
	iidiv = number.IFloorDiv                         //整数向下取整
	fidiv = number.FFloorDiv                         //浮点向下取整
	band  = func(a, b int64) int64 { return a & b }  //位运算 &
	bor   = func(a, b int64) int64 { return a | b }  //位运算 |
	bxor  = func(a, b int64) int64 { return a ^ b }  //位运算 ^
	shl   = number.ShiftLeft                         //左位移 <<
	shr   = number.ShiftRight                        //右位移 >>
	iunm  = func(a, _ int64) int64 { return -a }     //整数自取反 -x
	funm  = func(a, _ float64) float64 { return -a } //浮点数自取反 -x
	bnot  = func(a, _ int64) int64 { return ^a }     //条件 取反

)

//哪种类型 的操作
type operator struct {
	integerFunc func(int64, int64) int64
	floatFunc   func(float64, float64) float64
}

//保持 和 const.go 运算符 的一致
var operators = []operator{
	operator{iadd, fadd},
	operator{isub, fsub},
	operator{imul, fmul},
	operator{imod, fmod},
	operator{nil, pow},
	operator{nil, div},
	operator{iidiv, fidiv},
	operator{band, nil},
	operator{bor, nil},
	operator{bxor, nil},
	operator{shl, nil},
	operator{shr, nil},
	operator{iunm, funm},
	operator{bnot, nil},
}

func (self *luaState) Arith(op ArithOp) {
	var a, b luaValue
	b = self.stack.pop()
	if op != LUA_OPUNM && op != LUA_OPBNOT {
		a = self.stack.pop()
	} else {
		a = b
	}

	operator := operators[op]
	if result := _arith(a, b, operator); result != nil {
		self.stack.push(result)
	} else {
		panic("arithmetic error!") //这个操作必定有返回值
	}
}

func _arith(a, b luaValue, op operator) luaValue {
	if op.floatFunc == nil { //按位计算的操作
		if x, ok := convertToInteger(a); ok {
			if y, ok := convertToInteger(a); ok {
				return op.integerFunc(x, y)
			}
		}
	} else {
		if op.integerFunc != nil { //add,sub,mul,mod,idiv,unm
			if x, ok := a.(int64); ok {
				if y, ok := b.(int64); ok {
					return op.integerFunc(x, y)
				}
			}
		}

		if x, ok := convertToFloat(a); ok {
			if y, ok := convertToFloat(b); ok {
				return op.floatFunc(x, y)
			}
		}
	}

	//否则返回nil  错误
	return nil
}
