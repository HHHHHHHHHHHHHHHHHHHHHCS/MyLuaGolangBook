package parser

import . "LuaGo/compiler/ast"
import . "LuaGo/compiler/lexer"

//编译
func Parse(chunk, chunkName string) *Block {
	lexer := NewLexer(chunk, chunkName)
	block := parseBlock(lexer)
	lexer.NextTokenOfKind(TOKEN_EOF)
	return block
}
