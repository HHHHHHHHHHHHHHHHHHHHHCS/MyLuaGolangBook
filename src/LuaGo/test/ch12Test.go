package test

import (
	. "LuaGo/api"
	"LuaGo/state"
	"fmt"
	"io/ioutil"
)

type CH12Test struct {
}

func (test *CH12Test) DoTest() {
	data, err := ioutil.ReadFile("src/CH00_Luac/CH12Test.luac")
	if err != nil {
		panic(err)
	}
	ls := state.New()
	ls.Register("print", test.print)
	ls.Register("next", test.next)
	ls.Register("pairs", test.pairs)
	ls.Register("ipairs", test.iPairs)

	ls.Load(data, "chunk", "b")
	ls.Call(0, 0)
}

func (test CH12Test) print(ls BasicAPI) int {
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

func (test CH12Test) next(ls BasicAPI) int {
	//参数1 是表   参数2 是上一个键   如果2没有则为nil
	ls.SetTop(2)
	if ls.Next(1) {
		//返回2个参数 key value
		return 2
	} else {
		//返回1个参数 nil
		ls.PushNil()
		return 1
	}
}

func (test CH12Test) pairs(ls BasicAPI) int {
	//3个参数      第三个是头指针是nil
	ls.PushGoFunction(test.next)
	ls.PushValue(1)
	ls.PushNil()
	return 3
}

func (test CH12Test) _iPairsAux(ls BasicAPI) int {
	i := ls.ToInteger(2) + 1
	ls.PushInteger(i)
	if ls.GetI(1, i) == LUA_TNIL {
		//返回1个参数 nil
		return 1
	} else {
		//返回2个参数
		return 2
	}
}

func (test CH12Test) iPairs(ls BasicAPI) int {
	ls.PushGoFunction(test._iPairsAux)
	ls.PushValue(1)
	ls.PushInteger(0)
	return 3
}
