package test

import (
	. "LuaGo/api"
	"LuaGo/state"
	"fmt"
)

type CH05Test struct {
}

func (test *CH05Test) DoTest() {
	ls := state.New(20, nil)
	ls.PushInteger(1)
	ls.PushString("2.0")
	ls.PushString("3.0")
	ls.PushNumber(4.0)
	test.printStack(ls)

	ls.Arith(LUA_OPADD)
	test.printStack(ls)
	ls.Arith(LUA_OPBNOT)
	test.printStack(ls)
	ls.Len(2)
	test.printStack(ls)
	ls.Concat(3)
	test.printStack(ls)
	ls.PushBoolean(ls.Compare(1, 2, LUA_OPEQ))
	test.printStack(ls)
}

func (test *CH05Test) printStack(ls LuaState) {
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
