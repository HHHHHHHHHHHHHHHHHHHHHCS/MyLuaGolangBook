package main

import (
	"LuaGo/binchunk"
	"fmt"

	//"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}
		proto := binchunk.Undump(data)
		println(proto)
	}
}

func list(f *binchunk.Prototype) {

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
		fmt.Printf("\t%d\t%s\n", i+1)
	}
}

func constantToString(k interface{}) string {
	//https://studygolang.com/articles/2644
	switch k.(type) {
	case nil:
		return ""
	case bool:
		return fmt.Sprintf("%t", k)

	}
}
