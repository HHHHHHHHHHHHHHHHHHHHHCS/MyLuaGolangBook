package test

import (
	. "LuaGo/compiler/lexer"
	"fmt"
	"io/ioutil"
)

type CH14Test struct {
}

func (test *CH14Test) DoTest() {
	path := "src/CH00_Luac/CH14Test.lua"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	test.testLexer(string(data), path)
}

func (test *CH14Test) testLexer(chunk, chunkName string) {
	lexer := NewLexer(chunk, chunkName)
	for {
		line, kind, token := lexer.NextToken()
		fmt.Printf("[%2d][%-10s] %s\n",
			line, test.kindToCategory(kind), token)
		if kind == TOKEN_EOF {
			break
		}
	}
}

func (test *CH14Test) kindToCategory(kind int) string {
	switch {
	case kind < TOKEN_SEP_SEMI:
		return "other"
	case kind <= TOKEN_SEP_RCURLY:
		return "separator"
	case kind <= TOKEN_OP_NOT:
		return "operator"
	case kind <= TOKEN_KW_WHILE:
		return "keyword"
	case kind == TOKEN_IDENTIFIER:
		return "identifier"
	case kind == TOKEN_NUMBER:
		return "number"
	case kind == TOKEN_STRING:
		return "string"
	default:
		return "other"
	}
}
