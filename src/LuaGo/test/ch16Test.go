package test

import (
	"LuaGo/compiler/parser"
	"encoding/json"
	"io/ioutil"
)

type CH16Test struct {
}

func (test *CH16Test) DoTest() {
	path := "src/CH00_Luac/CH16Test.lua"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	test.testParser(string(data), "MyChunk")
}

func (test *CH16Test) testParser(chunk, chunkName string) {
	ast := parser.Parse(chunk, chunkName)
	b, err := json.Marshal(ast)
	if err != nil {
		panic(err)
	}
	println(string(b))
}
