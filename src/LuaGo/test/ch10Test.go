package test


import (
	. "LuaGo/api"
	"LuaGo/state"
	"fmt"
	"io/ioutil"
)

type CH10Test struct {
}

func (test *CH10Test) DoTest() {
	data, err := ioutil.ReadFile("src/CH00_Luac/CH10Test.luac")
	if err != nil {
		panic(err)
	}
	ls := state.New()
	ls.Register("print", test.print)
	ls.Load(data, "chunk", "b")
	ls.Call(0, 0)
}

func (test CH10Test) print(ls LuaState) int {
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
