package state

type luaState struct {
	stack *luaStack
}

var DefaultStackSize int = 20

//Lua栈初始容量
func New(stackSize int) *luaState {
	return &luaState{
		stack: newLuaStack(stackSize),
	}
}
