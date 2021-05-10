package stdlib

import (
	. "LuaGo/api"
	"os"
	"strings"
)

//key in the registry for table of loaded moodules
const LUA_LOADED_TABLE = "_LOADED"

//key in the registry for table of preloaded loaders
const LUA_PRELOAD_TABLE = "_PRELOAD"

const (
	LUA_DIRSEP    = string(os.PathSeparator) //'\\'
	LUA_PATH_SEP  = ";"
	LUA_PATH_MARK = "?"
	LUA_EXEC_DIR  = "!"
	LUA_IGMARK    = "-"
)

var llFuncs = map[string]GoFunction{
	"require": pkgRequire,
}

var pkgFuncs = map[string]GoFunction{
	"searchpath": pkgSearchPath,
}

func OpenPackageLib(ls LuaState) int {
	ls.NewLib(pkgFuncs)
	createSearchsTable(ls)
	//set paths
	ls.PushString("./?.lua;./?/init.lua")
	ls.SetField(-2, "path")
	//store config information
	ls.PushString(LUA_DIRSEP + "\n" + LUA_PATH_SEP + "\n" +
		LUA_PATH_MARK + "\n" + LUA_EXEC_DIR + "\n" + LUA_IGMARK + "\n")
	ls.SetField(-2, "config")
	//set field 'loaded'
	ls.GetSubTable(LUA_REGISTRYINDEX, LUA_LOADED_TABLE)
	ls.SetField(-2, "loaded")
	//set field 'preload'
	ls.GetSubTable(LUA_REGISTRYINDEX, LUA_PRELOAD_TABLE)
	ls.SetField(-2, "preload")
	ls.PushGlobalTable()
	ls.PushValue(-2)        //set 'package' as upvalue for next lib
	ls.SetFuncs(llFuncs, 1) //open lib into global table
	ls.Pop(1)               //pop global table
	return 1                //return 'package' table
}

func createSearchsTable(ls LuaState) {
	searchers := []GoFunction{
		preloadSearcher,
		luaSearcher,
	}
	ls.CreateTable(len(searchers), 0)
	for idx, searcher := range searchers {
		ls.PushValue(-2)
		ls.PushGoClosure(searcher, 1)
		ls.RawSetI(-2, int64(idx+1))
	}
	ls.SetField(-2, "searchers")
}

func preloadSearcher(ls LuaState) int {
	name := ls.CheckString(1)
	ls.GetField(LUA_REGISTRYINDEX, "_PRELOAD")
	if ls.GetField(-1, name) == LUA_TNIL { //not found
		ls.PushString("\n\tno field package.preload['" + name + "']")
	}
	return 1
}

func luaSearcher(ls LuaState) int {
	name := ls.CheckString(1)
	ls.GetField(LuaUpvalueIndex(1), "path")
	path, ok := ls.ToStringX(-1)
	if !ok {
		ls.Error2("'package.path' must be a string!")
	}

	filename, errMsg := _searchPath(name, path, ".", LUA_DIRSEP)
	if errMsg != "" {
		ls.PushString(errMsg)
		return 1
	}

	if ls.LoadFile(filename) == LUA_OK {
		ls.PushString(filename)
		return 2
	} else {
		return ls.Error2("error loading module '%s' from file '%s':\n\t%s",
			ls.CheckString(1), filename, ls.CheckString(-1))
	}
}

func _searchPath(name, path, sep, dirSep string) (filename, errMsg string) {
	if sep != "" {
		name = strings.Replace(name, sep, dirSep, -1)
	}

	for _, filename := range strings.Split(path, LUA_PATH_SEP) {
		filename = strings.Replace(filename, LUA_PATH_MARK, name, -1)
		if _, err := os.Stat(filename); !os.IsNotExist(err) {
			return filename, ""
		}
		errMsg += "\n\tno file '" + filename + "'"
	}
	return "", errMsg
}

func pkgSearchPath(ls LuaState) int {
	name := ls.CheckString(1)
	path := ls.CheckString(2)
	//把. 转换成 \\  地址用
	sep := ls.OptString(3, ".")
	rep := ls.OptString(4, LUA_DIRSEP)
	if filename, errMsg := _searchPath(name, path, sep, rep); errMsg == "" {
		ls.PushString(filename)
		return 1

	} else {
		ls.PushNil()
		ls.PushString(errMsg)
		return 2
	}
}

func pkgRequire(ls LuaState) int {
	name:=ls.CheckString(1)
	//todo:
}
