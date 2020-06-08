package main

import (
	"LuaGo/binchunk"
	"fmt"
	"io/ioutil"
)

func main() {
	data, err := ioutil.ReadFile("src/CH00_Luac/luac.out")
	if err != nil {
		panic(err)
	}
	proto := binchunk.Undump(data)
	println(proto)
	/*
		if len(os.Args) > 1 {
			data, err := ioutil.ReadFile(os.Args[1])
			if err != nil {
				panic(err)
			}
			proto := binchunk.Undump(data)
			println(proto)
		}
	*/
}

func list(f *binchunk.Prototype) {
	printHeader(f)
	printCode(f)
	printDetail(f)
	for _, p := range f.Protos {
		list(p)
	}
}

func printHeader(f *binchunk.Prototype) {
	funcType := "main"
	if f.LineDefined > 0 {
		funcType = "function"
	}

	varargFlag := ""

	if f.IsVararg > 0 {
		varargFlag = "+"
	}

	fmt.Printf("\n%s <%s:%d,%d (%dinstructions)\n", funcType,
		f.Source, f.LineDefined, f.LastLineDefined, len(f.Code))

	fmt.Printf("%d%s params , %d slots , %d upvalues, ",
		f.NumParams, varargFlag, f.MaxStackSize, len(f.Upvalues))

	fmt.Printf("%d locals, %d constants, %d functions\n",
		len(f.LocVars), len(f.Constants), len(f.Protos))
}

func printCode(f *binchunk.Prototype) {
	for pc, c := range f.Code {
		line := "-"
		if len(f.LineInfo) > 0 {
			line = fmt.Sprintf("%d", f.LineInfo[pc])
		}
		//打印 指令的序号 行号 十六进制
		fmt.Printf("\t%d\t[%s]\t0x%08X\n", pc+1, line, c)
	}
}

func printDetail(f *binchunk.Prototype) {
	fmt.Printf("constants (%d):\n", len(f.Constants))
	for i, k := range f.Constants {
		fmt.Printf("\t%d\t%s\n", i+1, constantToString(k))
	}

	fmt.Printf("locals (%d):\n", len(f.LocVars))
	for i, locvar := range f.LocVars {
		fmt.Printf("\t%d\t%s\t%d\t%d\n",
			i, locvar.VarName, locvar.StartPC+1, locvar.EndPC+1)
	}

	fmt.Printf("upvalues (%d):\n", len(f.Upvalues))
	for i, upval := range f.Upvalues {
		fmt.Printf("\t%d\t%s\t%d\t%d\n",
			i, upvalName(f, i), upval.Instack, upval.Idx)
	}
}

func constantToString(k interface{}) string {
	//https://studygolang.com/articles/2644
	switch k.(type) {
	case nil:
		return ""
	case bool:
		return fmt.Sprintf("%t", k)
	case float64:
		return fmt.Sprintf("%g", k)
	case int64:
		return fmt.Sprintf("%d", k)
	case string:
		return fmt.Sprintf("%q", k)
	default:
		return "?"
	}
}

func upvalName(f *binchunk.Prototype, idx int) string {
	if len(f.UpValueNames) > 0 {
		return f.UpValueNames[idx]
	}

	return "-"
}
