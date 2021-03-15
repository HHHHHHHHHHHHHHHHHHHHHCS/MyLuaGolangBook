package parser

import . "LuaGo/compiler/ast"
import . "LuaGo/compiler/lexer"

func _finishPrefixExp(lexer *Lexer, exp Exp) Exp {
	for {
		switch lexer.LookAhead() {
		case TOKEN_SEP_LBRACK:
			lexer.NextToken() //`[`
			keyExp := parseExp(lexer)
			lexer.NextTokenOfKind(TOKEN_SEP_RBRACK) //`]`
			exp = &TableAccessExp{LastLine: lexer.Line(), PrefixExp: exp, KeyExp: keyExp}
		case TOKEN_SEP_DOT:
			lexer.NextToken()                    //`.`
			line, name := lexer.NextIdentifier() //name
			keyExp := &StringExp{Line: line, Str: name}
			exp = &TableAccessExp{LastLine: line, PrefixExp: exp, KeyExp: keyExp}
		case TOKEN_SEP_COLON, TOKEN_SEP_LABEL,
			TOKEN_SEP_LCURLY, TOKEN_STRING:
			exp = _finishFuncCallExp(lexer, exp) //[`:` name] args
		default:
			return exp
		}
		return exp
	}
}

func parseParensExp(lexer *Lexer) Exp {
	lexer.NextTokenOfKind(TOKEN_SEP_LPAREN) //`(`
	exp := parseExp(lexer)                  //exp
	lexer.NextTokenOfKind(TOKEN_SEP_RPAREN) //`)`
	switch exp.(type) {
	case *VarargExp, *FuncCallExp, *NameExp, *TableAccessExp:
		return &ParensExp{Exp: exp}
	}
	return exp
}

func parsePrefixExp(lexer *Lexer) Exp {
	var exp Exp
	if lexer.LookAhead() == TOKEN_IDENTIFIER {
		line, name := lexer.NextIdentifier() //name
		exp = &NameExp{Line: line, Name: name}
	} else {
		exp = parseParensExp(lexer)
	}
	return _finishPrefixExp(lexer, exp)
}
