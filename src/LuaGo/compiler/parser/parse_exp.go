package parser

import . "LuaGo/compiler/ast"
import . "LuaGo/compiler/lexer"

// 把 , 的表达式连接起来
func parseExpList(lexer *Lexer) []Exp {
	exps := make([]Exp, 0, 4)
	exps = append(exps, parseExp(lexer))
	for lexer.LookAhead() == TOKEN_SEP_COMMA { //','
		lexer.NextToken()
		exps = append(exps, parseExp(lexer))
	}
}
