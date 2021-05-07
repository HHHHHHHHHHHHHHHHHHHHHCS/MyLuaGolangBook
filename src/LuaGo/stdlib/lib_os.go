package stdlib

//import "C" 需要下载 MinGW
// https://sourceforge.net/projects/mingw-w64/
// 然后配置环境变量

// #include <time.h>
import "C"

import (
	. "LuaGo/api"
	"os"
	"time"
)

var sysLib = map[string]GoFunction{
	"clock":     osClock,
	"difftime":  osDiffTime,
	"time":      osTime,
	"date":      osDate,
	"remove":    osRemove,
	"rename":    osRename,
	"tmpname":   osTmpName,
	"getenv":    osGetEnv,
	"execute":   osExecute,
	"exit":      osExit,
	"setlocale": osSetLocale,
}

func OpenOSLib(ls LuaState) int {
	ls.NewLib(sysLib)
	return 1
}

//os.clock()
func osClock(ls LuaState) int {
	c := float64(C.clock()) / float64(C.CLOCKS_PER_SEC)
	ls.PushNumber(c)
	return 1
}

//os.difftime(t2,t1)
func osDiffTime(ls LuaState) int {
	t2 := ls.CheckInteger(1)
	t1 := ls.CheckInteger(2)
	ls.PushInteger(t2 - t1)
	return 1
}

//os.time([table])
func osTime(ls LuaState) int {
	if ls.IsNoneOrNil(1) { //without args?
		t := time.Now().Unix()
		ls.PushInteger(t)
	} else {
		ls.CheckType(1, LUA_TTABLE)
		sec := _getField(ls, "sec", 0)
		min := _getField(ls, "min", 0)
		hour := _getField(ls, "hour", 12)
		day := _getField(ls, "day", -1)
		month := _getField(ls, "month", -1)
		year := _getField(ls, "year", -1)

		//time.Local 当前区
		t := time.Date(year, time.Month(month), day,
			hour, min, sec, 0, time.Local).Unix()
		ls.PushInteger(t)
	}
	return 1
}

//dft : get field value if is nil then default value
func _getField(ls LuaState, key string, dft int64) int {
	t := ls.GetField(-1, key)
	res, isNum := ls.ToIntegerX(-1)
	if !isNum {
		if t != LUA_TNIL {
			return ls.Error2("field '%s' is not an integer", key)
		} else if dft < 0 {
			return ls.Error2("field '%s' missing in date table", key)
		}
		res = dft
	}
	ls.Pop(1)
	return int(res)
}

//os.date([format [, time]]
func osDate(ls LuaState) int {
	format := ls.OptString(1, "%c")
	var t time.Time
	if ls.IsInteger(2) {
		t = time.Unix(ls.ToInteger(2), 0)
	} else {
		t = time.Now()
	}

	if format != "" && format[0] == '!' { //utc?
		format = format[1:]
		t = t.In(time.UTC)
	}

	if format == "*t" {
		ls.CreateTable(0, 9)
		_setField(ls, "sec", t.Second())
		_setField(ls, "min", t.Minute())
		_setField(ls, "hour", t.Hour())
		_setField(ls, "day", t.Day())
		_setField(ls, "month", int(t.Month()))
		_setField(ls, "year", t.Year())
		_setField(ls, "wday", int(t.Weekday())+1)
		_setField(ls, "yday", t.YearDay())
	} else if format == "%c" {
		ls.PushString(t.Format(time.ANSIC))
	} else {
		ls.PushString(format)
	}

	return 1
}

func _setField(ls LuaState, key string, value int) {
	ls.PushInteger(int64(value))
	ls.SetField(-2, key)
}

//os.remove(filename)
func osRemove(ls LuaState) int {
	filename := ls.CheckString(1)
	if err := os.Remove(filename); err != nil {
		ls.PushNil()
		ls.PushString(err.Error())
		return 2
	} else {
		ls.PushBoolean(true)
		return 1
	}
}

//os.rename(oldname,newname)
func osRename(ls LuaState) int {
	oldName := ls.CheckString(1)
	newName := ls.CheckString(2)
	if err := os.Rename(oldName, newName); err != nil {
		ls.PushNil()
		ls.PushString(err.Error())
		return 2
	} else {
		ls.PushBoolean(true)
		return 1
	}
}

//os.tmpname()
func osTmpName(ls LuaState) int {
	panic("todo: osTmpName!")
}

//os.getenv(varname)
//返回环境变量  如java_home配置什么的
func osGetEnv(ls LuaState) int {
	key := ls.CheckString(1)
	if env := os.Getenv(key); env != "" {
		ls.PushString(env)
	} else {
		ls.PushNil()
	}
	return 1
}

//os.execute([command])
func osExecute(ls LuaState) int {
	panic("todo: osExecute!")
}

//os.exit([code [, close]])
func osExit(ls LuaState) int {
	if ls.IsBoolean(1) {
		if ls.ToBoolean(1) {
			os.Exit(0)
		} else {
			os.Exit(1) //todo
		}
	} else {
		code := ls.OptInteger(1, 1)
		os.Exit(int(code))
	}
	if ls.ToBoolean(2) {
		//ls.Close()
	}

	return 0
}

//os.setlocale(locale [, category])
func osSetLocale(ls LuaState) int {
	panic("todo: osSetLocale!")
}
