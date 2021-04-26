package test

import (
	. "LuaGo/api"
	"LuaGo/state"
	"fmt"
	"io/ioutil"
)

type CH13Test struct {
}

func (test *CH13Test) DoTest() {
	data, err := ioutil.ReadFile("src/CH00_Luac/CH13Test.luac")
	if err != nil {
		panic(err)
	}
	ls := state.New()
	ls.Register("print", test.print)
	ls.Register("error", test.error)
	ls.Register("pcall", test.pCall)

	ls.Load(data, "chunk", "b")
	ls.Call(0, 0)
}

func (test CH13Test) print(ls BasicAPI) int {
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

func (test CH13Test) error(ls BasicAPI) int {
	//默认只接受一个参数 错误对象在栈顶
	return ls.Error()
}

func (test CH13Test) pCall(ls BasicAPI) int {
	nArgs := ls.GetTop() - 1 //第一个是栈堆  后面一位才是参数个数
	status := ls.PCall(nArgs, -1, 0)
	//因为是后面push进来的   所以这里要翻转下
	ls.PushBoolean(status == LUA_OK)
	ls.Insert(1)
	return ls.GetTop()
}
