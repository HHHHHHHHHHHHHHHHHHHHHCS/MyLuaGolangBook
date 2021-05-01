package stdlib

import (
	. "LuaGo/api"
	"fmt"
	"strconv"
	"strings"
)

var baseFuncs = map[string]GoFunction{
	"print":        basePrint,
	"assert":       baseAssert,
	"error":        baseError,
	"select":       baseSelect,
	"ipairs":       baseIPairs,
	"pairs":        basePairs,
	"next":         baseNext,
	"load":         baseLoad,
	"loadfile":     baseLoadFile,
	"dofile":       baseDoFile,
	"pcall":        basePCall,
	"xpcall":       baseXPCall,
	"getmetatable": baseGetMetatable,
	"setmetatable": baseSetMetatable,
	"rawequal":     baseRawEqual,
	"rawlen":       baseRawLen,
	"rawget":       baseRawGet,
	"rawset":       baseRawSet,
	"type":         baseType,
	"tostring":     baseToString,
	"tonumber":     baseToNumber,
	/* placeholders */
	"_G":       nil,
	"_VERSION": nil,
}

//注册基础的API func
func OpenBaseLib(ls LuaState) int {
	//open lib into global table
	ls.PushGlobalTable()
	ls.SetFuncs(baseFuncs, 0)
	//set global _G
	ls.PushValue(-1)
	ls.SetField(-2, "_G")
	//set global _Version
	ls.PushString("Lua 5.3")
	ls.SetField(-2, "_VERSION")
	return 1
}

//print
func basePrint(ls LuaState) int {
	n := ls.GetTop() //number of arg
	ls.GetGlobal("tostring")
	for i := 1; i <= n; i++ {
		ls.PushValue(-1)
		ls.PushValue(i)
		ls.Call(1, 1)
		s, ok := ls.ToStringX(-1)
		if !ok {
			return ls.Error2("'tostring' must return a string to 'print'")
		}
		if i > 1 {
			fmt.Print("\t")
		}
		fmt.Print(s)
		ls.Pop(1)
	}
	fmt.Println()
	return 0
}

//assert
func baseAssert(ls LuaState) int {
	if ls.ToBoolean(1) { //判断bool
		return ls.GetTop() //如果true 直接返回bool
	} else {
		ls.CheckAny(1)                     //如果是false
		ls.Remove(1)                       //删除顶部
		ls.PushString("assertion failed!") //错误log
		ls.SetTop(1)                       //把错误信息写到顶部
		return baseError(ls)               //然后输出错误信息
	}
}

//error
func baseError(ls LuaState) int {
	level := int(ls.OptInteger(2, 1))
	ls.SetTop(1)
	if ls.Type(1) == LUA_TSTRING && level > 0 {
		//ls.Where(level)
		//ls.PushValue(1)
		//ls.Concat(2)
	}
	return ls.Error()
}

//select
//select(1, 'a', 'b', 'c') ---> 'a'
//select(#, 'a', 'b', 'c') ---> 'a' , 'b' , 'c'
func baseSelect(ls LuaState) int {
	n := int64(ls.GetTop())
	if ls.Type(1) == LUA_TSTRING && ls.CheckString(1) == "#" {
		//#返回参数个数
		ls.PushInteger(n - 1)
		return 1
	} else {
		i := ls.CheckInteger(1)
		if i < 0 {
			i = n + i
		} else if i > n {
			i = n
		}
		ls.ArgCheck(1 <= i, 1, "index out of range")
		//把适当数量的参数作为返回值 就好了
		return int(n - i)
	}
}

//iparis
func baseIPairs(ls LuaState) int {
	ls.CheckAny(1)
	ls.PushGoFunction(iPairsAux) //iteration function
	ls.PushValue(1)              //state
	ls.PushInteger(0)            //initial value
	return 3
}

func iPairsAux(ls LuaState) int {
	i := ls.CheckInteger(2) + 1
	ls.PushInteger(i)
	if ls.GetI(1, i) == LUA_TNIL {
		return 1
	} else {
		return 2
	}
}

//pairs
func basePairs(ls LuaState) int {
	ls.CheckAny(1)
	if ls.GetMetafield(1, "__pairs") == LUA_TNIL {
		ls.PushGoFunction(baseNext)
		ls.PushValue(1)
		ls.PushNil()
	} else {
		ls.PushValue(1) //use  self
		ls.Call(1, 3)
	}
	return 3
}

func baseNext(ls LuaState) int {
	ls.CheckType(1, LUA_TTABLE)
	ls.SetTop(2)
	if ls.Next(1) {
		return 2
	} else {
		ls.PushNil()
		return 1
	}
}

//load
func baseLoad(ls LuaState) int {
	var status int
	chunk, isStr := ls.ToStringX(1)
	mode := ls.OptString(3, "bt")
	env := 0
	if !ls.IsNone(4) {
		env = 4
	}
	if isStr {
		chunkname := ls.OptString(2, chunk)
		status = ls.Load([]byte(chunk), chunkname, mode)
	} else {
		panic("loading from a reader function")
	}
	return loadAux(ls, status, env)
}

func loadAux(ls LuaState, status, envIdx int) int {
	if status == LUA_OK {
		if envIdx != 0 {
			panic("todo!")
		}
		return 1
	} else {
		ls.PushNil()
		ls.Insert(-2)
		return 2
	}
}

//loadfile
func baseLoadFile(ls LuaState) int {
	fname := ls.OptString(1, "")
	mode := ls.OptString(1, "bt")
	env := 0
	if !ls.IsNone(3) {
		env = 3
	}
	status := ls.LoadFileX(fname, mode)
	return loadAux(ls, status, env)
}

//dofile
func baseDoFile(ls LuaState) int {
	fname := ls.OptString(1, "bt")
	ls.SetTop(1)
	if ls.LoadFile(fname) != LUA_OK {
		return ls.Error()
	}
	ls.Call(0, LUA_MULTRET)
	return ls.GetTop() - 1
}

//pcall
func basePCall(ls LuaState) int {
	nArgs := ls.GetTop() - 1
	status := ls.PCall(nArgs, -1, 0)
	ls.PushBoolean(status == LUA_OK)
	ls.Insert(1)
	return ls.GetTop()
}

func baseXPCall(ls LuaState) int {
	panic("todo!")
}

//getmetatable
func baseGetMetatable(ls LuaState) int {
	ls.CheckAny(1)
	if !ls.GetMetatable(1) {
		ls.PushNil()
		return 1
	}
	ls.GetMetafield(1, "__metatable")
	return 1
}

//setmetatable
func baseSetMetatable(ls LuaState) int {
	t := ls.Type(2)
	ls.CheckType(1, LUA_TTABLE)
	ls.ArgCheck(t == LUA_TNIL || t == LUA_TTABLE, 2,
		"nil or table expected")
	if ls.GetMetafield(1, "__metatable") != LUA_TNIL {
		return ls.Error2("cannot change a protected metatable")
	}
	ls.SetTop(2)
	ls.SetMetatable(1)
	return 1
}

//a==b
func baseRawEqual(ls LuaState) int {
	ls.CheckAny(1)
	ls.CheckAny(2)
	ls.PushBoolean(ls.RawEqual(1, 2))
	return 1
}

//table or string length
func baseRawLen(ls LuaState) int {
	t := ls.Type(1)
	ls.ArgCheck(t == LUA_TTABLE || t == LUA_TSTRING, 1,
		"table or string expected")
	ls.PushInteger(int64(ls.RawLen(1)))
	return 1
}

//get table[xxx]
func baseRawGet(ls LuaState) int {
	ls.CheckType(1, LUA_TTABLE)
	ls.CheckAny(2)
	ls.SetTop(2)
	ls.RawGet(1)
	return 1
}

//set table[xxx] = xxx
func baseRawSet(ls LuaState) int {
	ls.CheckType(1, LUA_TTABLE)
	ls.CheckAny(2)
	ls.CheckAny(3)
	ls.SetTop(3)
	ls.RawSet(1)
	return 1
}

//type name
func baseType(ls LuaState) int {
	t := ls.Type(1)
	ls.ArgCheck(t != LUA_TNONE, 1, "value expected")
	ls.PushString(ls.TypeName(t))
	return 1
}

//tostring
func baseToString(ls LuaState) int {
	ls.CheckAny(1)
	ls.ToString2(1)
	return 1
}

//进制转换
func baseToNumber(ls LuaState) int {
	//
	if ls.IsNoneOrNil(2) {
		ls.CheckAny(1)
		if ls.Type(1) == LUA_TNUMBER {
			ls.SetTop(1)
			return 1
		} else {
			if s, ok := ls.ToStringX(1); ok {
				if ls.StringToNumber(s) {
					return 1
				}
			}
		}
	} else {
		//进制转换
		ls.CheckType(1, LUA_TSTRING)
		s := strings.TrimSpace(ls.ToString(1))
		base := int(ls.CheckInteger(2))
		ls.ArgCheck(2 <= base && base <= 36, 2, "base out of range")
		if n, err := strconv.ParseInt(s, base, 64); err == nil {
			ls.PushInteger(n)
			return 1
		}
	}
	ls.PushNil()
	return 1
}
