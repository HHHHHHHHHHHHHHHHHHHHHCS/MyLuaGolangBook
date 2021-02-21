package ast

/*
stat ::=  ‘;’ |
	 varlist ‘=’ explist |
	 functioncall |
	 label |
	 break |
	 goto Name |
	 do block end |
	 while exp do block end |
	 repeat block until exp |
	 if exp then block {elseif exp then block} [else block] end |
	 for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end |
	 for namelist in explist do block end |
	 function funcname funcbody |
	 local function Name funcbody |
	 local namelist [‘=’ explist]
*/
type Stat interface {
}

//`;`
//空语句没有任何意义
type EmptyStat struct{}

//break
//break 因为要跳转行 所以要记录行
type BreakStat struct {
	Line int
}

//`::`Name`::`
//goto的name用
type LabelStat struct {
	Name string
}

//goto Name
//goto到哪个name
type GotoStat struct {
	Name string
}

//do block end
//执行某个代码块
type DoStat struct {
	Block *Block
}

//functional
//调用函数 可以是语句也可以是表达式
type FuncCallStat = FuncCallExp

//while循环
//while exp do block end
type WhileStat struct {
	Exp   Exp
	Block *Block
}

//repeat循环
//repeat block until exp
type RepeatStat struct {
	Block *Block
	Exp   Exp
}

//if语句
//if exp then block {elseif exp then block} [else block] end
//如果把最后的else 看做 elseif true 则可以缩写为
//if exp then block {elseif exp then block} [elseif true then block] end
//if exp then block {elseif exp then block} end
//索引0是 if then 其余都是 elseif
type IfStat struct {
	Exps   []Exp
	Blocks []*Block
}

//for循环
//for Name '=' exp ',' [',' exp] do block end
type ForNumStat struct {
	LineOfFor int
	LineOfDo  int
	VarName   string
	InitExp   Exp
	LimitExp  Exp
	StepExp   Exp
	Block     *Block
}

//for in 循环
//for nameList in explist do block end
// namelist ::= Name {‘,’ Name}
// explist ::= exp {‘,’ exp}
type ForInStat struct {
	LineOfDo int
	NameList []string
	ExpList  []Exp
	Block    *Block
}

//local变量申明
//要把末尾行号记录下来 以供生成阶段使用
//local namelist ['=' explist]
//namelist ::= Name {',' Name}
//explist ::= exp{',' exp}
type LocalVarDeclStat struct {
	LastLine int
	NameList []string
	ExpList  []Exp
}

//赋值语句
//varlist '=' explist
//varlist ::= var {',' var}
//var :: Name | prefixexp '[' exp ']' | prefixexp '.' Name
//explist ::= exp {',' exp}
type AssignStat struct {
	LastLine int
	VarList  []Exp
	ExpList  []Exp
}

//非局部函数
//function funcName funcBody
//	=> function t.a.b.c:f(param) body end
//funcName ::= Name {'.' Name} [':' Name]
//	=> function t.a.b.c.f(self,param) body end
//funcBody ::= '(' [parlist] ')' block end
//	=> t.a.b.c.f = function(self, params) body end
type LocalFuncDefStat struct {
	Name string
	Exp  *FuncDefExp
}
