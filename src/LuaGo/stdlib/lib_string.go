package stdlib

import . "Luago/api"

var strLib = map[string]GoFunction{
	"len":      strLen,
	"rep":      strRep,
	"reverse":  strReverse,
	"lower":    strLower,
	"upper":    strUpper,
	"sub":      strSub,
	"byte":     strByte,
	"char":     strChar,
	"dump":     strDump,
	"format":   strFormat,
	"packsize": strPackSize,
	"pack":     strPack,
	"unpack":   strUnpack,
	"find":     strFind,
	"match":    strMatch,
	"gsub":     strGsub,
	"gmatch":   strGmatch,
}

func OpenStringLib(ls LuaState) int{
	ls.NewLib(strLib)
	createMetatable(ls)
	return 1
}

func createMetatable(ls LuaState)  {
	ls.CreateTable(0, 1)
	ls.PushString("dummy")
	ls.PushValue(-2)
	ls.SetMetatable(-2)
	ls.Pop(1)
	ls.PushValue(-2)
	ls.SetField(-2, "__index")
	ls.Pop(1)
}
