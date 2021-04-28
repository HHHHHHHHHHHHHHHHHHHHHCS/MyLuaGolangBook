package test

import (
	"LuaGo/api"
	"LuaGo/state"
	"fmt"
	"io/ioutil"
)

type CH17Test struct {
}

func (test *CH17Test) DoTest() {
	path := "src/CH00_Luac/CH17Test.lua"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	ls := state.New()
	ls.Register("print", test.__print)
	ls.Register("getmetatable", test.__getMetatable)
	ls.Register("setmetatable", test.__setMetatable)
	ls.Register("next", test.__next)
	ls.Register("pairs", test.__pairs)
	ls.Register("ipairs", test.__iPairs)
	ls.Register("error", test.__error)
	ls.Register("pcall", test.__pCall)
	ls.Load(data, "MyChunk", "bt")
	ls.Call(0, 0)
}


func (test *CH17Test) __print(ls api.BasicAPI) int {
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

func (test *CH17Test) __getMetatable(ls api.BasicAPI) int {
	if !ls.GetMetatable(1) {
		ls.PushNil()
	}
	return 1
}

func (test *CH17Test) __setMetatable(ls api.BasicAPI) int {
	ls.SetMetatable(1)
	return 1
}

func (test *CH17Test) __next(ls api.BasicAPI) int {
	ls.SetTop(2) /* create a 2nd argument if there isn't one */
	if ls.Next(1) {
		return 2
	} else {
		ls.PushNil()
		return 1
	}
}

func (test *CH17Test) __pairs(ls api.BasicAPI) int {
	ls.PushGoFunction(test.__next) /* will return generator, */
	ls.PushValue(1)           /* state, */
	ls.PushNil()
	return 3
}

func (test *CH17Test) __iPairs(ls api.BasicAPI) int {
	ls.PushGoFunction(test.___iPairsAux) /* iteration function */
	ls.PushValue(1)                 /* state */
	ls.PushInteger(0)               /* initial value */
	return 3
}

func (test *CH17Test) ___iPairsAux(ls  api.BasicAPI) int {
	i := ls.ToInteger(2) + 1
	ls.PushInteger(i)
	if ls.GetI(1, i) ==  api.LUA_TNIL {
		return 1
	} else {
		return 2
	}
}

func (test *CH17Test) __error(ls api.BasicAPI) int {
	return ls.Error()
}

func (test *CH17Test) __pCall(ls api.BasicAPI) int {
	nArgs := ls.GetTop() - 1
	status := ls.PCall(nArgs, -1, 0)
	ls.PushBoolean(status == api.LUA_OK)
	ls.Insert(1)
	return ls.GetTop()
}
