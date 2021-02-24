package parser

import (
	. "LuaGo/compiler/ast"
)
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

func parseIfStat(lexer *Lexer) *IfStat {
	exps := make([]Exp, 0, 4)
	blocks := make([]*Block, 0, 4)

	lexer.NextTokenOfKind(TOKEN_KW_IF)         //if
	exps = append(exps, parseExp(lexer))       //exp
	lexer.NextTokenOfKind(TOKEN_KW_THEN)       //then
	blocks = append(blocks, parseBlock(lexer)) //block

	for lexer.LookAhead() == TOKEN_KW_ELSEIF { //{
		lexer.NextToken()                          //elseif
		exps = append(exps, parseExp(lexer))       //exp
		lexer.NextTokenOfKind(TOKEN_KW_THEN)       //then
		blocks = append(blocks, parseBlock(lexer)) //block
	} //}

	//else block => elseif true then block
	if lexer.LookAhead() == TOKEN_KW_ELSE { //{
		lexer.NextToken()                                 //else
		exps = append(exps, &TrueExp{Line: lexer.Line()}) //true
		blocks = append(blocks, parseBlock(lexer))        //block
	} //}

	lexer.NextTokenOfKind(TOKEN_KW_END)
	return &IfStat{Exps: exps, Blocks: blocks} //end
}

func parseForStat(lexer *Lexer) Stat {
	lineOfFor, _ = lexer.NextTokenOfKind(TOKEN_KW_FOR)
	_, name := lexer.NextIdentifier()
	//有 = 号  默认为 for ; ; do
	//这里是偷懒的做法 要思考补齐
	if lexer.LookAhead() == TOKEN_OP_ASSIGN {
		return _finishForNumStat(lexer, lineOfFor, name)
	} else {
		return _finishForInStat(lexer, name)
	}
}

func _finishForNumStat(lexer *Lexer, lineOfFor int, varName string) *ForNumStat {
	lexer.NextTokenOfKind(TOKEN_OP_ASSIGN) //for name `=`
	initExp := parseExp(lexer)             //exp
	lexer.NextTokenOfKind(TOKEN_SEP_COMMA) //`,`
	limitExp := parseExp(lexer)

	var stepExp Exp
	if lexer.LookAhead() == TOKEN_SEP_COMMA {
		lexer.NextToken()         //`,`
		stepExp = parseExp(lexer) //exp
	} else {
		stepExp = &IntegerExp{Line: lexer.Line(), Val: 1} //默认+1
	}

	lineOfDo, _ := lexer.NextTokenOfKind(TOKEN_KW_DO) //do
	block := parseBlock(lexer)                        //block
	lexer.NextTokenOfKind(TOKEN_KW_END)               //end

	return &ForNumStat{LineOfFor: lineOfFor, LineOfDo: lineOfDo,
		VarName: varName, InitExp: initExp, LimitExp: limitExp, StepExp: stepExp, Block: block}
}

func _finishNameList(lexer *Lexer, name0 string) []string {
	names := []string{name0}
	for lexer.LookAhead() == TOKEN_SEP_COMMA {
		lexer.NextToken()                 //`,`
		_, name := lexer.NextIdentifier() //name
		names = append(names, name)
	}
	return names
}

func _finishForInStat(lexer *Lexer, name0 string) *ForInStat {
	nameList := _finishNameList(lexer, name0) //for nameList
	lexer.NextTokenOfKind(TOKEN_KW_IN)
	expList := parseExpList(lexer)
	lineOfDo, _ := lexer.NextTokenOfKind(TOKEN_KW_DO)
	block := parseBlock(lexer)
	lexer.NextTokenOfKind(TOKEN_KW_END)
	return &ForInStat{LineOfDo: lineOfDo, NameList: nameList, ExpList: expList, Block: block}
}
