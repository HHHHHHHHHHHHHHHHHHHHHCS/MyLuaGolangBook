package parser

import . "LuaGo/compiler/ast"
import . "LuaGo/compiler/lexer"

func parseStat(lexer *Lexer) Stat {
	switch lexer.LookAhead() {
	case TOKEN_SEP_SEMI:
		return parseEmptyStat(lexer)
	case TOKEN_KW_BREAK:
		return parseBreakStat(lexer)
	case TOKEN_SEP_LABEL:
		return parseLabelStat(lexer)
	case TOKEN_KW_GOTO:
		return parseGotoStat(lexer)
	case TOKEN_KW_DO:
		return parseDoStat(lexer)
	case TOKEN_KW_WHILE:
		return parseWhileStat(lexer)
	case TOKEN_KW_REPEAT:
		return parseRepeatStat(lexer)
	case TOKEN_KW_IF:
		return parseIfStat(lexer)
	case TOKEN_KW_FOR:
		return parseForStat(lexer)
	case TOKEN_KW_FUNCTION:
		return parseFuncDefStat(lexer)
	case TOKEN_KW_LOCAL:
		return parseLocalAssignOrFuncDefStat(lexer)
	default:
		return parseAssignOrFuncCallStat(lexer)
	}
}

func parseEmptyStat(lexer *Lexer) *EmptyStat {
	lexer.NextTokenOfKind(TOKEN_SEP_SEMI) //';'
	return &EmptyStat{}
}

func parseBreakStat(lexer *Lexer) *BreakStat {
	lexer.NextTokenOfKind(TOKEN_KW_BREAK) //break
	return &BreakStat{Line: lexer.Line()}
}

func parseLabelStat(lexer *Lexer) *LabelStat {
	lexer.NextTokenOfKind(TOKEN_SEP_LABEL) //`::`
	_, name := lexer.NextIdentifier()      //Name
	lexer.NextTokenOfKind(TOKEN_SEP_LABEL) //`::`
	return &LabelStat{Name: name}
}

func parseGotoStat(lexer *Lexer) *GotoStat {
	lexer.NextTokenOfKind(TOKEN_KW_GOTO) //goto
	_, name := lexer.NextIdentifier()    //Name
	return &GotoStat{Name: name}
}

func parseDoStat(lexer *Lexer) *DoStat {
	lexer.NextTokenOfKind(TOKEN_KW_DO)  //do
	block := parseBlock(lexer)          //block
	lexer.NextTokenOfKind(TOKEN_KW_END) //end
	return &DoStat{Block: block}
}

func parseWhileStat(lexer *Lexer) *WhileStat {
	lexer.NextTokenOfKind(TOKEN_KW_WHILE) //while
	exp := parseExp(lexer)                //exp
	block := parseBlock(lexer)            //block
	lexer.NextTokenOfKind(TOKEN_KW_END)   //end
	return &WhileStat{Exp: exp, Block: block}
}

func parseRepeatStat(lexer *Lexer) *RepeatStat {
	lexer.NextTokenOfKind(TOKEN_KW_REPEAT) //repeat
	block := parseBlock(lexer)             //block
	lexer.NextTokenOfKind(TOKEN_KW_UNTIL)  //until
	exp := parseExp(lexer)                 //exp
	return &RepeatStat{Block: block, Exp: exp}
}
