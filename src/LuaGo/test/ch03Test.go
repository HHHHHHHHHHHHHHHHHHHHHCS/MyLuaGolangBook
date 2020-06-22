package test

//.表示导入以后，该包的函数和变量不需要再直接写入包名称。
//_的作用就更特殊。当导入一个包的时候，该包的init和其他函数都会被导入；但不是所有函数都需要. "_"符号可以只导入init，而不需要导入其他函数。
import (
	"LuaGo/binchunk"
	. "LuaGo/vm"
	"fmt"
	"io/ioutil"
)

type CH03Test struct {
}

func (test *CH03Test) DoTest() {
	data, err := ioutil.ReadFile("src/CH00_Luac/luac.out")
	if err != nil {
		panic(err)
	}
	proto := binchunk.Undump(data)
	test.list(proto)
	/*
		if len(os.Args) > 1 {
			data, err := ioutil.ReadFile(os.Args[1])
			if err != nil {
				panic(err)
			}
			proto := binchunk.Undump(data)
			list(proto)
		}
	*/
}

func (test *CH03Test) list(f *binchunk.Prototype) {
	test.printHeader(f)
	test.printCode(f)
	test.printDetail(f)
	for _, p := range f.Protos {
		test.list(p)
	}
}

func (test *CH03Test) printHeader(f *binchunk.Prototype) {
	funcType := "main"
	if f.LineDefined > 0 {
		funcType = "function"
	}

	varargFlag := ""

	if f.IsVararg > 0 {
		varargFlag = "+"
	}

	fmt.Printf("\n%s <%s:%d,%d> (%d instructions)\n",
		funcType, f.Source, f.LineDefined, f.LastLineDefined, len(f.Code))

	fmt.Printf("%d%s params, %d slots, %d upvalues, ",
		f.NumParams, varargFlag, f.MaxStackSize, len(f.Upvalues))

	fmt.Printf("%d locals, %d constants, %d functions\n",
		len(f.LocVars), len(f.Constants), len(f.Protos))
}

func (test *CH03Test) printCode(f *binchunk.Prototype) {
	for pc, c := range f.Code {
		line := "-"
		if len(f.LineInfo) > 0 {
			line = fmt.Sprintf("%d", f.LineInfo[pc])
		}
		i := Instruction(c)
		//打印 指令的序号 行号 十六进制
		fmt.Printf("\t%d\t[%s]\t0x%08X\n", pc+1, line, c)
		printOperands(i)
		fmt.Printf("\n")
	}
}

func (test *CH03Test) printDetail(f *binchunk.Prototype) {
	fmt.Printf("constants (%d):\n", len(f.Constants))
	for i, k := range f.Constants {
		fmt.Printf("\t%d\t%s\n", i+1, test.constantToString(k))
	}

	fmt.Printf("locals (%d):\n", len(f.LocVars))
	for i, locVar := range f.LocVars {
		fmt.Printf("\t%d\t%s\t%d\t%d\n",
			i, locVar.VarName, locVar.StartPC+1, locVar.EndPC+1)
	}

	fmt.Printf("upvalues (%d):\n", len(f.Upvalues))
	for i, upval := range f.Upvalues {
		fmt.Printf("\t%d\t%s\t%d\t%d\n",
			i, test.upvalName(f, i), upval.Instack, upval.Idx)
	}
}

func (test *CH03Test) constantToString(k interface{}) string {
	//https://studygolang.com/articles/2644
	switch k.(type) {
	case nil:
		return "nil"
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

func (test *CH03Test) upvalName(f *binchunk.Prototype, idx int) string {
	if len(f.UpvalueNames) > 0 {
		return f.UpvalueNames[idx]
	}
	return "-"
}

func printOperands(i Instruction) {
	switch i.OpMode() {
	case IABC:
		a, b, c := i.ABC()
		fmt.Printf("%d", a)
		if i.BMode() != OpArgN {
			if b > 0xFF {
				fmt.Printf("%d", -1-b&0xFF)
			} else {
				fmt.Printf(" %d", b)
			}
		}
		if i.CMode() != OpArgN {
			if c > 0xFF {
				fmt.Printf(" %d", -1-c&0xFF)
			} else {
				fmt.Printf(" %d", c)
			}
		}
	}
}
