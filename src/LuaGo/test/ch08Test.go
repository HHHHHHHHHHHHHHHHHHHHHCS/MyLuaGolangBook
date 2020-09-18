package test

import (
	"LuaGo/state"
	"io/ioutil"
)

type CH08Test struct {
}

func (test *CH08Test) DoTest() {
	data, err := ioutil.ReadFile("src/CH00_Luac/CH08Test.luac")
	if err != nil {
		panic(err)
	}
	ls := state.New()
	ls.Load(data, "src/CH00_Luac/CH08Test.luac", "b")
	ls.Call(0, 0)
}

//TODO:
