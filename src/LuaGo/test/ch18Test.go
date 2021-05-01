package test

import (
	"LuaGo/state"
)

type CH18Test struct {
}

func (test *CH18Test) DoTest() {
	path := "src/CH00_Luac/CH18Test.lua"
	//data, err := ioutil.ReadFile(path)
	//if err != nil {
	//	panic(err)
	//}

	ls := state.New()
	ls.OpenLibs()
	ls.LoadFile(path)
	ls.Call(0, -1)
}
