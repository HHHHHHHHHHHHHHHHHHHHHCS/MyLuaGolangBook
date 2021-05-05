package stdlib

import (
	. "LuaGo/api"
	"unicode/utf8"
)

const UTF8PATT = "[\x00-\x7F\xC2-\xF4][\x80-\xBF]*"

const MAX_UNICODE = 0x10FFFF

var utf8Lib = map[string]GoFunction{
	"len":       utfLen,
	"offset":    utfByteOffset,
	"codepoint": utfCodePoint,
	"char":      utfChar,
	"codes":     utfIterCodes,
	/* placeholders */
	"charpattern": nil,
}

func OpenUTF8Lib(ls LuaState) int {
	ls.NewLib(utf8Lib)
	ls.PushString(UTF8PATT)
	ls.SetField(-2, "charpattern")
	return 1
}

//utf8.len(s,[, i [, j]])
func utfLen(ls LuaState) int {
	s := ls.CheckString(1)
	sLen := len(s)
	i := posRelat(ls.OptInteger(2, 1), sLen)
	j := posRelat(ls.OptInteger(3, -1), sLen)
	ls.ArgCheck(1 <= i && i <= sLen+1, 2,
		"initial position out of string")
	ls.ArgCheck(j <= sLen, 3,
		"final position out of string")

	if i > j {
		ls.PushInteger(0)
	} else {
		n := utf8.RuneCountInString(s[i-1 : j])
		ls.PushInteger(int64(n))
	}

	return 1
}

//utf8.offset(s,n [, i])
func utfByteOffset(ls LuaState) int {
	s := ls.CheckString(1)
	sLen := len(s)
	n := ls.CheckInteger(2)
	i := 1
	if n < 0 {
		i = sLen + 1
	}
	i = posRelat(ls.OptInteger(3,int64()))
	//todo:

}
