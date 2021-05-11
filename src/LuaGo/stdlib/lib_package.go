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

var pkgFuncs = map[string]GoFunction{
	"searchpath": pkgSearchPath,
	/* placeholders */
	"preload":   nil,
	"cpath":     nil,
	"path":      nil,
	"searchers": nil,
	"loaded":    nil,
}

var llFuncs = map[string]GoFunction{
	"require": pkgRequire,
}

func OpenPackageLib(ls LuaState) int {
	ls.NewLib(pkgFuncs)
	createSearchersTable(ls)
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

func createSearchersTable(ls LuaState) {
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
		ls.Error2("'package.path' must be a string")
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

func pkgRequire(ls LuaState) int {
	name := ls.CheckString(1)
	ls.SetTop(1) //LOADED table will be at index 2
	ls.GetField(LUA_REGISTRYINDEX, LUA_LOADED_TABLE)
	ls.GetField(2, name) //LOADED[name]
	if ls.ToBoolean(-1) {
		return 1 //package is already loaded
	}
	//else
	ls.Pop(1) // remove 'getfield' result
	_findLoader(ls, name)
	ls.PushString(name) //pass name arg
	ls.Insert(-2)       //name is 1st arg (before search data)
	ls.Call(2, 1)       //run loader to load module
	if !ls.IsNil(-1) {  //LOADED[name] = return value
		ls.SetField(2, name)
	}
	if ls.GetField(2, name) == LUA_TNIL { //module set no value
		ls.PushBoolean(true) //use true as result
		ls.PushValue(-1)     //extra copy to be returned
		ls.SetField(2, name) //LOADED[name] = true
	}
	return 1
}

func _findLoader(ls LuaState, name string) {
	if ls.GetField(LuaUpvalueIndex(1), "searchers") != LUA_TTABLE {
		ls.Error2("'package.searchers' must be a table")
	}

	errMsg := "module '" + name + "' not found:"

	for i := int64(1); ; i++ {
		if ls.RawGetI(3, i) == LUA_TNIL { //no more searchers?
			ls.Pop(1)         //remove nil
			ls.Error2(errMsg) //create error msg
		}

		ls.PushString(name)
		ls.Call(1, 2)          //call
		if ls.IsFunction(-2) { //find a loader
			return
		} else if ls.IsString(-2) { //searcher returned error msg
			ls.Pop(1)
			errMsg += ls.CheckString(-1)
		} else { //remove both returns
			ls.Pop(2)
		}
	}
}
