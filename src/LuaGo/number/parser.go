package number

import "strconv"

//TODO:NEED CHANGE

//转换成整数
func ParseInteger(str string) (int64, bool) {
	//string  进制  位数(比如int32 int64
	i, err := strconv.ParseInt(str, 10, 64)
	return i, err == nil
}

//转换为float
func ParseFloat(str string) (float64, bool) {
	f, err := strconv.ParseFloat(str, 64)
	return f, err == nil
}
