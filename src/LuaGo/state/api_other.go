package state

import (
	. "LuaGo/api"
	"fmt"
)

func (ls *luaState) RegisterOtherAPI() {
	ls.Register("print", __print)
	ls.Register("getmetatable", __getMetatable)
	ls.Register("setmetatable", __setMetatable)
	ls.Register("next", __next)
	ls.Register("pairs", __pairs)
	ls.Register("ipairs", __iPairs)
	ls.Register("error", __error)
	ls.Register("pcall", __pCall)
}

func __print(ls BasicAPI) int {
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

func __getMetatable(ls BasicAPI) int {
	if !ls.GetMetatable(1) {
		ls.PushNil()
	}
	return 1
}

func __setMetatable(ls BasicAPI) int {
	ls.SetMetatable(1)
	return 1
}

func __next(ls BasicAPI) int {
	ls.SetTop(2) /* create a 2nd argument if there isn't one */
	if ls.Next(1) {
		return 2
	} else {
		ls.PushNil()
		return 1
	}
}

func __pairs(ls BasicAPI) int {
	ls.PushGoFunction(__next) /* will return generator, */
	ls.PushValue(1)              /* state, */
	ls.PushNil()
	return 3
}

func __iPairs(ls BasicAPI) int {
	ls.PushGoFunction(___iPairsAux) /* iteration function */
	ls.PushValue(1)                    /* state */
	ls.PushInteger(0)                  /* initial value */
	return 3
}

func ___iPairsAux(ls BasicAPI) int {
	i := ls.ToInteger(2) + 1
	ls.PushInteger(i)
	if ls.GetI(1, i) == LUA_TNIL {
		return 1
	} else {
		return 2
	}
}

func __error(ls BasicAPI) int {
	return ls.Error()
}

func __pCall(ls BasicAPI) int {
	nArgs := ls.GetTop() - 1
	status := ls.PCall(nArgs, -1, 0)
	ls.PushBoolean(status == LUA_OK)
	ls.Insert(1)
	return ls.GetTop()
}
