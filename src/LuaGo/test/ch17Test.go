package test

import (
	"LuaGo/state"
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
	ls.RegisterOtherAPI()
	ls.Load(data, "MyChunk", "bt")
	ls.Call(0, 0)
}
