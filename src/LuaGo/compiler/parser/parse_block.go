package parser

import . "LuaGo/compiler/ast"
import . "LuaGo/compiler/lexer"

func parseBlock(lexer *Lexer) *Block {
	return &Block{
		Stats:    parseStats(lexer),   //行为语句
		RetExps:  parseRetExps(lexer), //return表达式
		LastLine: lexer.Line(),
	}
}

func _isReturnOrBlockEnd(tokenKind int) bool {
	switch tokenKind {
	case TOKEN_KW_RETURN, TOKEN_EOF, TOKEN_KW_END,
		TOKEN_KW_ELSE, TOKEN_KW_ELSEIF, TOKEN_KW_UNTIL:
		return true
	}
	return false
}

func parseRetExps(lexer *Lexer) []Exp {
	if lexer.LookAhead() != TOKEN_KW_RETURN {
		return nil
	}

	lexer.NextToken()
	switch lexer.LookAhead() {
	case TOKEN_EOF, TOKEN_KW_END,
		TOKEN_KW_ELSE, TOKEN_KW_IF, TOKEN_KW_UNTIL:
		return []Exp{}
	case TOKEN_SEP_SEMI:
		lexer.NextToken()
		return []Exp{}
	default:
		exps := parseExpList(lexer)
		if lexer.LookAhead() == TOKEN_SEP_SEMI { //';'
			lexer.NextToken()
		}
		return exps
	}
}

func parseStats(lexer *Lexer) []Stat {
	stats := make([]Stat, 0, 8)
	for !_isReturnOrBlockEnd(lexer.LookAhead()) {
		stat := parseStat(lexer)
		if _, ok := stat.(*EmptyStat); !ok {
			stats = append(stats, stat)
		}
	}
	return stats
}
