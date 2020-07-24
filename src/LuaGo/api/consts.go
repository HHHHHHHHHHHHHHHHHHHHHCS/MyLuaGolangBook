package api

const (
	LUA_TNONE = iota - 1 //-1 无效索引值
	LUA_TNIL
	LUA_TBOOLEAN
	LUA_TLIGHTUSERDATA
	LUA_TNUMBER
	LUA_TSTRING
	LUA_TTABLE
	LUA_TFUNCTION
	LUA_TUSERDATA
	LUA_TTHREAD
)

//运算符
const (
	LUA_OPADD  = iota // +
	LUA_OPSUB         // -
	LUA_OPMUL         // *
	LUA_OPMOD         // %
	LUA_OPPOW         // ^
	LUA_OPDIV         // /
	LUA_OPIDIV        // // 整除
	LUA_OPBAND        // &
	LUA_OPBOR         // |
	LUA_OPBXOR        // ~.
	LUA_OPSHL         // <<
	LUA_OPSHR         // >>
	LUA_OPNUM         // - 自负数
	LUA_OPBNOT        // ~ 条件 取反

)

//比较符
const (
	LUA_OPEQ = iota // ==
	LUA_OPLT        // <
	LUA_OPLE        //<=
)
