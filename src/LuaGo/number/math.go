package number

import "math"

//go没有直接的乘法运算符
//go整除只适用于整数 , 并且直接截断 , 并非向下取整(负数)
//go取模不能直接映射的

//向下取整
func IFloorDiv(a, b int64) int64 {
	if a > 0 && b > 0 || a < 0 && b < 0 || a%b == 0 {
		return a / b
	} else {
		return a/b - 1
	}
}

//向下取整
func FFloorDiv(a, b float64) float64 {
	return math.Floor(a / b)
}

//取模 整除模拟
func IMod(a, b int64) int64 {
	return a - IFloorDiv(a, b)*b
}

//取模 整除模拟
func FMod(a, b float64) float64 {
	return a - math.Floor(a/b)*b
}

//左位移  右边操作数只能是无符号整数
//如果n<0 则代表反向位移 (因为GO 只支持 正方向 移动)
func ShiftLeft(a, n int64) int64 {
	if n >= 0 {
		return a << uint64(n)
	} else {
		return ShiftRight(a, -n)
	}
}

//右位移  我们需要的是无符号位移
//所以先转换成无符号 处理好了再转换为有符号的
func ShiftRight(a, n int64) int64 {
	if n >= 0 {
		return int64(uint64(a) >> uint64(n))
	} else {
		return ShiftLeft(a, -n)
	}
}

//浮点数转换成整数  去除小数点 且没有超出范围 则表示成功
func FloatToInteger(f float64) (int64, bool) {
	i := int64(f)
	return i, float64(i) == f
}
