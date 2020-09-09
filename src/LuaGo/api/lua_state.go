package api

type LuaType = int
type ArithOp = int
type CompareOp = int

//栈基础操作函数
type LuaState interface {
	//基础栈操作
	GetTop() int
	AbsIndex(idx int) int
	CheckStack(n int) bool
	Pop(n int)
	Copy(fromIdx, toIdx int)
	PushValue(idx int)
	Replace(idx int)
	Insert(idx int)
	Remove(idx int)
	Rotate(idx, n int)
	SetTop(idx int)
	//栈访问 stack->go
	TypeName(tp LuaType) string
	Type(idx int) LuaType
	IsNone(idx int) bool
	IsNil(idx int) bool
	IsNoneOrNil(idx int) bool
	IsBoolean(idx int) bool
	IsInteger(idx int) bool
	IsNumber(idx int) bool
	IsString(idx int) bool
	IsTable(idx int) bool
	IsThread(idx int) bool
	IsFunction(idx int) bool
	ToBoolean(idx int) bool
	ToInteger(idx int) int64
	ToIntegerX(idx int) (int64, bool)
	ToNumber(idx int) float64
	ToNumberX(idx int) (float64, bool)
	ToString(idx int) string
	ToStringX(idx int) (string, bool)
	//压入栈 go->stacks
	PushNil()
	PushBoolean(b bool)
	PushInteger(n int64)
	PushNumber(n float64)
	PushString(s string)
	//执行算数和按位运算
	Arith(op ArithOp)
	Compare(idx1, idx2 int, op CompareOp) bool
	//其他方法
	Len(idx int)
	Concat(n int)
	//Table
	CreateTable(nArr, nRec int)
	GetTable(idx int) LuaType
	GetField(idx int, k string) LuaType
	GetI(idx int, i int64) LuaType
	SetTable(idx int)
	SetField(idx int, k string)
	SetI(idx int, n int64)
	//Function
	Load(chunk []byte, chunkName, mode string) int
	Call(nArgs, nResults int)
}
