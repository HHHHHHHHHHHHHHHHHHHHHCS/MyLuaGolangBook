package test

import (
	"LuaGo/state"
	"fmt"
	"io/ioutil"
)
import . "LuaGo/api"

type CH09Test struct {
}

func (test *CH09Test) DoTest() {
	data, err := ioutil.ReadFile("src/CH00_Luac/CH09Test.luac")
	if err != nil {
		panic(err)
	}
	ls := state.New()
	ls.Register("print", test.print)
	ls.Load(data, "chunk", "b")
	ls.Call(0, 0)
}

func (test CH09Test) print(ls LuaState) int {
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
