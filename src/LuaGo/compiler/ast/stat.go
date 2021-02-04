package ast

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
//todo的name用
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

//while exp do block end
type whileStat struct {
	Exp   Exp
	Block *Block
}

//repeat block until exp
type RepeatStat struct {
	Block *Block
	Exp   Exp
}

//if exp then block {elseif exp then block} [else block] end
//如果把最后的else 看做 elseif true 则可以缩写为
//if exp then block {elseif exp then block} [elseif true then block] end
//if exp then block {elseif exp then block} end
//索引0是 if then 其余都是 elseif
type IfStat struct {
	Exps   []Exp
	Blocks []*Blocks
}

//for Name '=' exp ',' [',' exp] do block end
type ForNumStat struct {
	LineOfFor int
	LineOfDo  int
	VarName   string
	InitExp   Exp
	StepExp   Exp
	Block     *Block
}

//for nameList in explist do block end
type ForInStat struct {
	LineOfDo int
	NameList []string
	ExpList  []Exp
	Block    *Block
}

//local namelist ['=' explist]
//namelist ::= Name {',' Name}
//explist ::= exp{',' exp}
type LocalVarDeclStat struct {
	LastLine int
	NameList []string
	ExpList  []Exp
}
