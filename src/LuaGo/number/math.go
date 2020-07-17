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
