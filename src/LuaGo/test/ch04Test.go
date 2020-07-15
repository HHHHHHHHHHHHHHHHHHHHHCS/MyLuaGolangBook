package test

import (
	"fmt"
	. "LuaGo/api"
	"LuaGo/state"
)

type CH04Test struct {
}

func (test *CH04Test) DoTest() {
	ls:=state.New()
	ls.PushBoolean(true)
	printStack(ls)

}

func printStack(ls LuaState){
	ls.GetTop()
}
