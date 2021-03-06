package test

import (
	. "LuaGo/api"
	"LuaGo/state"
	"fmt"
	"io/ioutil"
)

type CH11Test struct {
}

func (test *CH11Test) DoTest() {
	data, err := ioutil.ReadFile("src/CH00_Luac/CH11Test.luac")
	if err != nil {
		panic(err)
	}
	ls := state.New()
	ls.Register("print", test.print)
	ls.Register("getmetatable", test.getMetatable)
	ls.Register("setmetatable", test.setMetatable)

	ls.Load(data, "chunk", "b")
	ls.Call(0, 0)
}

func (test CH11Test) print(ls LuaState) int {
	nArgs := ls.GetTop()
	for i := 1; i <= nArgs; i++ {
		if ls.IsBoolean(i) {
			fmt.Printf("%t", ls.ToBoolean(i))
		} else if ls.IsString(i) {
			fmt.Print(ls.ToString(i))
		} else {
			fmt.Print(ls.TypeName(ls.Type(i)))
		}
		if i < nArgs {
			fmt.Print("\t")
		}
	}
	fmt.Println()
	return 0
}

func (test CH11Test) getMetatable(ls LuaState) int {
	if !ls.GetMetatable(1) {
		ls.PushNil()
	}
	return 1
}

func (test CH11Test) setMetatable(ls LuaState) int {
	ls.SetMetatable(1)
	return 1
}
