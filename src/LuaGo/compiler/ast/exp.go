package ast

type Exp interface{}

//简单的赋值表达式
//---------------------
type NilExp struct{ Line int }
type TrueExp struct{ Line int }
type FalseExp struct{ Line int }
type VarargExp struct{ Line int }
type IntegerExp struct {
	Line int
	Val  int64
}
type FloatExp struct {
	Line int
	Val  float64
}
type StringExp struct {
	Line int
	Str  string
}

//运算符表达式
//------------

//一元运算符
type UnopExp struct {
	Line int
	Op   int
	Exp  Exp
}

//二元运算符
type BinopExp struct {
	Line int
	Op   int
	Exp1 Exp
	Exp2 Exp
}

//拼接表达式
type ConcatExp struct {
	Line int
	Exps []Exp
}

//表的构造表达式
//tableconstructor ::= '{' [fieldlist] '}'
//fieldlist ::= field {fieldsep field}[fieldsep]
//field ::= '[' exp ']' '=' exp | Name '=' exp | exp
//fieldsep ::= ',' | ';'
type TableConstructorExp struct {
	Line     int //line of '{'
	LastLine int //line of '}'
	KeyExps  []Exp
	ValExps  []Exp
}

//方法表达式
type FuncDefExp struct {
	Line     int
	LastLine int
	ParList  []string
	IsVararg bool
	Block    *Block
}

type NameExp struct {
	Line int
	Name string
}

//圆括号表达式
type ParensExp struct {
	Exp Exp
}

//表访问表达式
type TableAccessExp struct {
	LastLine  int //line of ']'
	PrefixExp Exp
	KeyExp    Exp
}

//使用方法表达式
type FuncCallExp struct {
	Line      int // line of '('
	LastLine  int // line of ')'
	PrefixExp Exp
	NameExp   *StringExp
	Args      []Exp
}

