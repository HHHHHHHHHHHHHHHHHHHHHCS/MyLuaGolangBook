package compiler

import (
	"LuaGo/binchunk"
	"LuaGo/compiler/codegen"
	"LuaGo/compiler/parser"
)

func Compile(chunk, chunkName string) *binchunk.Prototype {
	ast := parser.Parse(chunk, chunkName)
	return codegen.GenProto(ast)
}
