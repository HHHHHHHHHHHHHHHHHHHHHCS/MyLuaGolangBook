package stdlib

import (
	. "LuaGo/api"
	"sort"
	"strings"
)

const MAX_LEN = 1000000

const (
	TAB_R  = 1             //read
	TAB_W  = 2             //write
	TAB_L  = 4             //length
	TAB_RW = TAB_R | TAB_W //read|write
)

var tabFuncs = map[string]GoFunction{
	"move":   tabMove,
	"insert": tabInsert,
	"remove": tabRemove,
	"sort":   tabSort,
	"concat": tabConcat,
	"pack":   tabPack,
	"unpack": tabUnpack,
}

func OpenTableLib(ls LuaState) int {
	ls.NewLib(tabFuncs)
	return 1
}

// table.move (a1, f, e, t [,a2])
func tabMove(ls LuaState) int {
	f := ls.CheckInteger(2)
	e := ls.CheckInteger(3)
	t := ls.CheckInteger(4)
	tt := 1 //dest table
	if !ls.IsNoneOrNil(5) {
		tt = 5
	}
	_checkTab(ls, 1, TAB_R)
	_checkTab(ls, tt, TAB_W)
	if e >= f {
		var n, i int64
		ls.ArgCheck(f > 0 || e < LUA_MAXINTEGER+f, 3,
			"too many elements to move")
		n = e - f + 1
		ls.ArgCheck(t <= LUA_MAXINTEGER-n+1, 4,
			"destination wrap around")
		if t > e || t <= f || (tt != 1 && !ls.Compare(1, tt, LUA_OPEQ)) {
			for i = 0; i < n; i++ {
				ls.GetI(1, f+i)
				ls.SetI(tt, t+i)
			}
		} else {
			//倒序是如果过小可以先进行扩容
			for i = n - 1; i >= 0; i-- {
				ls.GetI(1, f+i)
				ls.SetI(tt, t+i)
			}
		}
	}
	ls.PushValue(tt)
	return 1
}

// table.insert (list, [pos,] value)
func tabInsert(ls LuaState) int {
	e := _auxGetN(ls, 1, TAB_RW) + 1 //first empty element
	var pos int64
	switch ls.GetTop() {
	case 2:
		pos = e
		break
	case 3:
		pos = ls.CheckInteger(2)
		ls.ArgCheck(1 <= pos && pos <= e, 2, "position out of bounds")
		for i := e; i > pos; i-- {
			ls.GetI(1, i-1)
			ls.SetI(1, i)
		}
		break
	default:
		return ls.Error2("wrong number of arguments to 'insert'")
	}
	ls.SetI(1, pos)
	return 0
}

// table.remove (list [, pos])
func tabRemove(ls LuaState) int {
	size := _auxGetN(ls, 1, TAB_RW)
	pos := ls.OptInteger(2, size)
	if pos != size {
		ls.ArgCheck(1 <= pos && pos <= size+1, 1, "position out of bounds")
	}
	//把value换到队伍尾巴再删除
	ls.GetI(1, pos)
	for ; pos < size; pos++ {
		ls.GetI(1, pos+1)
		ls.SetI(1, pos)
	}
	ls.PushNil()
	ls.SetI(1, pos)
	return 1
}

// table.concat (list [, sep [, i [, j]]])
func tabConcat(ls LuaState) int {
	tabLen := _auxGetN(ls, 1, TAB_R)
	sep := ls.OptString(2, "")
	i := ls.OptInteger(3, 1)
	j := ls.OptInteger(4, tabLen)

	if i > j {
		ls.PushString("")
		return 1
	}

	buf := make([]string, j-i+1)
	for k := i; k > 0 && k <= j; k++ {
		ls.GetI(1, k)
		if !ls.IsString(-1) {
			ls.Error2("invalid value (%s) at index %d in table for 'concat'", ls.TypeName2(-1), i)
		}
		buf[k-i] = ls.ToString(-1)
		ls.Pop(1)
	}
	ls.PushString(strings.Join(buf, sep))

	return 1
}

func _auxGetN(ls LuaState, n, w int) int64 {
	_checkTab(ls, n, w|TAB_L)
	return ls.Len2(n)
}

func _checkTab(ls LuaState, arg, what int) {
	if ls.Type(arg) != LUA_TTABLE {
		n := 1
		if ls.GetMetatable(arg) &&
			(what&TAB_R != 0 || _checkField(ls, "__index", &n)) &&
			(what&TAB_W != 0 || _checkField(ls, "__newindex", &n)) &&
			(what&TAB_L != 0 || _checkField(ls, "__len", &n)) {
			ls.Pop(n)
		} else {
			ls.CheckType(arg, LUA_TTABLE)
		}
	}
}

func _checkField(ls LuaState, key string, n *int) bool {
	ls.PushString(key)
	*n++
	return ls.RawGet(-*n) != LUA_TNIL
}

// table.pack (···)
//create new table   and   push/setField
func tabPack(ls LuaState) int {
	n := int64(ls.GetTop())
	ls.CreateTable(int(n), 1)
	ls.Insert(1)
	for i := n; i >= 1; i-- {
		ls.SetI(1, i)
	}
	ls.PushInteger(n)
	ls.SetField(1, "n")
	return 1
}

// table.unpack (list [, i [, j]])
// push i~j element to top
func tabUnpack(ls LuaState) int {
	i := ls.OptInteger(2, 1)
	e := ls.OptInteger(3, ls.Len2(1))
	if i > e {
		return 0
	}

	n := int(e - i + 1)
	if n <= 0 || n >= MAX_LEN || !ls.CheckStack(n) {
		return ls.Error2("too many results to unpack")
	}
	for ; i < e; i++ {
		ls.GetI(1, i)
	}
	ls.GetI(1, e)
	return n
}

// table.sort (list [, comp])
//go的sort 需要写比较方法 Len Less Swap 所以封装了一层
func tabSort(ls LuaState) int {
	w := wrapper{ls}
	ls.ArgCheck(w.Len() < MAX_LEN, 1, "array too big")
	sort.Sort(w)
	return 0
}

type wrapper struct {
	ls LuaState
}

func (self wrapper) Len() int {
	return int(self.ls.Len2(1))
}

func (self wrapper) Less(i, j int) bool {
	ls := self.ls
	if ls.IsFunction(2) {
		ls.PushValue(2)
		ls.GetI(1, int64(i+1))
		ls.GetI(1, int64(j+1))
		ls.Call(2, 1)
		b := ls.ToBoolean(-1)
		ls.Pop(1)
		return b
	} else {
		ls.GetI(1, int64(i+1))
		ls.GetI(1, int64(j+1))
		b := ls.Compare(-2, -1, LUA_OPLT)
		ls.Pop(2)
		return b
	}
}

//取出到栈顶  再重新set回去
func (self wrapper) Swap(i, j int) {
	ls := self.ls
	ls.GetI(1, int64(i+1))
	ls.GetI(1, int64(j+1))
	ls.SetI(1, int64(i+1))
	ls.SetI(1, int64(j+1))
}
