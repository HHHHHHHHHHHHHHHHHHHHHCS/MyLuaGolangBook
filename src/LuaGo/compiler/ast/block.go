package ast

type Block struct {
	LastLine int    //代码块末尾行号
	Stats    []Stat //语句 只能用于执行不能用于求值
	RetExps  []Exp  //表达式 只能用于求值 不能用于执行
}
