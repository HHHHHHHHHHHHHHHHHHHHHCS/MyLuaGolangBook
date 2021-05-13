package api

const LUA_MINSTACK = 20                         //默认最小栈长度
const LUAI_MAXSTACK = 1000000                   //最大栈索引 有正负
const LUA_REGISTRYINDEX = -LUAI_MAXSTACK - 1000 //LUA自己用的伪索引  返回组注册表用
const LUA_RIDX_GLOBALS int64 = 2                //全局环境注册表里面的索引
const LUA_RIDX_MAINTHREAD int64 = 1             //主线程索引
const LUA_MULTRET = -1                          //多返回值

const (
	LUA_MAXINTEGER = 1<<63 - 1
	LUA_MININTEGER = -1 << 63
)

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
	LUA_OPBXOR        // ~ 异或
	LUA_OPSHL         // <<
	LUA_OPSHR         // >>
	LUA_OPUNM         // - 自负数
	LUA_OPBNOT        // ~ 条件 取反

)

//比较符
const (
	LUA_OPEQ = iota // ==
	LUA_OPLT        // <
	LUA_OPLE        //<=
)

//异常
const (
	LUA_OK = iota //成功
	LUA_YIELD
	LUA_ERRRUN //PCall 用
	LUA_ERRSYNTAX
	LUA_ERRMEM
	LUA_ERRGCMM
	LUA_ERRERR
	LUA_ERRFILE
)
