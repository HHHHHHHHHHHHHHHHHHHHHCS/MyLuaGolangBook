package test

import (
	. "LuaGo/api"
	"LuaGo/state"
	"fmt"
)

type CH04Test struct {
}

func (test *CH04Test) DoTest() {
	ls := state.New()
	ls.PushBoolean(true)
	test.printStack(ls)
	ls.PushInteger(10)
	test.printStack(ls)
	ls.PushNil()
	test.printStack(ls)
	ls.PushString("hello")
	test.printStack(ls)
	ls.PushValue(-4)
	test.printStack(ls)
	ls.Replace(3)
	test.printStack(ls)
	ls.SetTop(6)
	test.printStack(ls)
	ls.Remove(-3)
	test.printStack(ls)
	ls.SetTop(-5)
	test.printStack(ls)
}

func (test *CH04Test) printStack(ls BasicAPI) {
	top := ls.GetTop()
	for i := 1; i <= top; i++ {
		t := ls.Type(i)
		switch t {
		case LUA_TBOOLEAN:
			fmt.Printf("[%t]", ls.ToBoolean(i))
		case LUA_TNUMBER:
			fmt.Printf("[%g]", ls.ToNumber(i))
		case LUA_TSTRING:
			fmt.Printf("[%q]", ls.ToString(i))
		default:
			fmt.Printf("[%s]", ls.TypeName(t))
		}
	}
	fmt.Println()
}
