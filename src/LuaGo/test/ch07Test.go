package test

import (
	. "LuaGo/api"
	"LuaGo/binchunk"
	"LuaGo/state"
	. "LuaGo/vm"
	"fmt"
	"io/ioutil"
)

type CH07Test struct {
}

func (test *CH07Test) DoTest() {
	/*
		if len(os.Args) > 1 {
			data, err := ioutil.ReadFile(os.Args[1])
			if err != nil {
				panic(err)
			}
			proto := binchunk.Undump(data)
			test.luaMain(proto)
		}
	*/
	data, err := ioutil.ReadFile("src/CH00_Luac/CH07Test.luac")
	if err != nil {
		panic(err)
	}
	proto := binchunk.Undump(data)
	test.luaMain(data, proto)
}

func (test *CH07Test) luaMain(data []byte, proto *binchunk.Prototype) {
	nRegs := int(proto.MaxStackSize)
	ls := state.New()
	ls.Load(data, "src/CH00_Luac/CH07Test.luac", "b")
	ls.SetTop(nRegs)
	for {
		pc := ls.PC()
		inst := Instruction(ls.Fetch())
		if inst.Opcode() != OP_RETURN {
			inst.Execute(ls)
			fmt.Printf("[%02d] %s", pc+1, inst.OpName())
			test.printStack(ls)
		} else {
			break
		}
	}
}

func (test *CH07Test) printStack(ls BasicAPI) {
	top := ls.GetTop()
	for i := 1; i <= top; i++ {
		t := ls.Type(i)
		switch t {
		case LUA_TBOOLEAN:
			fmt.Printf("[%t]", ls.ToBoolean(i))
		case LUA_TNUMBER:
			fmt.Printf("[%g]", ls.ToNumber(i))
		case LUA_TSTRING:
			fmt.Printf("[%q]", ls.ToString(i))
		default:
			fmt.Printf("[%s]", ls.TypeName(t))
		}
	}
	fmt.Println()
}
