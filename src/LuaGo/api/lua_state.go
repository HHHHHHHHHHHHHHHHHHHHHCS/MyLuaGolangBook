package api

type LuaType = int
type ArithOp = int
type CompareOp = int

type GoFunction func(LuaState) int

//伪索引 转换到 注册表索引
func LuaUpvalueIndex(i int) int {
	return LUA_REGISTRYINDEX - i
}

type LuaState interface {
	BasicAPI
	AuxLib
}

//栈基础操作函数
type BasicAPI interface {
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
	XMove(to LuaState, n int)
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
	IsGoFunction(idx int) bool
	ToBoolean(idx int) bool
	ToInteger(idx int) int64
	ToIntegerX(idx int) (int64, bool)
	ToNumber(idx int) float64
	ToNumberX(idx int) (float64, bool)
	ToString(idx int) string
	ToStringX(idx int) (string, bool)
	ToGoFunction(idx int) GoFunction
	ToThread(idx int) LuaState
	ToPointer(idx int) interface{}
	RawLen(idx int) uint
	//压入栈 go->stacks
	PushNil()
	PushBoolean(b bool)
	PushInteger(n int64)
	PushNumber(n float64)
	PushString(s string)
	PushFString(fmt string, a ...interface{})
	PushGoFunction(f GoFunction)
	PushGoClosure(f GoFunction, n int)
	PushGlobalTable()
	PushThread() bool
	//执行算数和按位运算
	Arith(op ArithOp)
	Compare(idx1, idx2 int, op CompareOp) bool
	RawEqual(idx1, idx2 int) bool
	//get lua->stack
	NewTable()
	CreateTable(nArr, nRec int)
	GetTable(idx int) LuaType
	GetField(idx int, k string) LuaType
	GetI(idx int, i int64) LuaType
	RawGet(idx int) LuaType
	RawGetI(idx int, i int64) LuaType
	GetMetatable(idx int) bool
	GetGlobal(name string) LuaType
	//set stack->lua
	SetTable(idx int)
	SetField(idx int, k string)
	SetI(idx int, i int64)
	RawSet(idx int)
	RawSetI(idx int, i int64)
	SetMetatable(idx int)
	SetGlobal(name string)
	Register(name string, f GoFunction)
	//Function
	Load(chunk []byte, chunkName, mode string) int
	Call(nArgs, nResults int)
	PCall(nArgs, nResults, msgh int) int
	//其他方法
	Len(idx int)
	Concat(n int)
	Next(idx int) bool
	//try catch
	Error() int
	StringToNumber(s string) bool
	//coroutine
	NewThread() LuaState
	Resume(from LuaState, nArgs int) int
	Yield(nResults int) int
	Status() int
	IsYieldable() bool
	GetStack() bool //debug
}
