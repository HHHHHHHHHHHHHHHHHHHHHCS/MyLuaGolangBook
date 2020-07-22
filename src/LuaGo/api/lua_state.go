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
	ToBoolean(idx int) bool
	ToInteger(idx int) int64
	ToIntegerX(idx int) (int64, bool)
	ToNumber(idx int) float64
	ToString(idx int) string
	ToStringX(idx int) (string, bool)
	//压入栈 go->stacks
	PushNil()
	PushInteger(n int64)
	PushNumber(n float64)
	PushString(n string)
}
