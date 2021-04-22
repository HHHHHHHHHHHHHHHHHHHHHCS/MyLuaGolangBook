package binchunk

const (
	LUA_SIGNATURE    = "\x1bLua"
	LUAC_VERSION     = 0x53
	LUAC_FORMAT      = 0
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	CINT_SIZE        = 4
	CSIZET_SIZE      = 8
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
)

//用空接口 代替 联合体(Union) 的效果 把不同的数据类型统一起来
const (
	TAG_NIL       = 0x00 //nil 不储存
	TAG_BOOLEAN   = 0x01 //bool 0,1
	TAG_NUMBER    = 0x03 //number Lua浮点数
	TAG_INTEGER   = 0x13 //integer Lua 整数
	TAG_SHORT_STR = 0x04 //string 短字符串
	TAG_LONG_STR  = 0x14 //string 长字符串
)

//二进制chunk
type binaryChunk struct {
	header                  //头部
	sizeUpvalues byte       //主函数upvalue的数量
	mainFunc     *Prototype //主函数原型
}

type header struct {
	//独特的签名 用于识别文件 ESCLua 的十六进制 0x1B4C7564  , ESC(Escape)
	signature [4]byte
	//版本号 大版本.小版本.发布号 5.4.5  => 5*16+4 => 83 发布号不统计
	version byte
	//格式号 默认0
	format byte
	//进一步校验  是不是lua文件或损坏 0x1993(lua发布年份) 0x0D(回车) 0x0A(换行) 0x1A(替换) 0x0A(新换行)
	luacData [6]byte
	//数据类型长度 不符合则拒绝加载
	//int 4位
	cintSize byte
	//size_t 8位
	sizetSize byte
	//lua虚拟机指令 4位
	instructionSize byte
	//lua整数 8位
	luaIntegerSize byte
	//lua浮点数 8位
	luaNumberSize byte
	//lua整数值 n(8)个字节存放 0x5678 目的检测二进制chunk的大小端方式
	luacInt int64
	//lua浮点数 n(8)个字节存放 370.5 目的检测二进制chunk的浮点数格式 如IEEE754
	luacNum float64
}

//函数原型
type Prototype struct {
	//main 函数  储存main函数的名字长度 + 符号 + 文件名字  符号@ 表示来自文件 #来自字符串
	Source string
	//起止行号 下面两个cint  普通函数起止行号大于0  主函数都是0
	LineDefined     uint32
	LastLineDefined uint32
	//固定参数个数 固定参数相对于变长参数而言的 主函数通常是0固定参数
	NumParams byte
	//是否有变长参数 0代表否 1代表是  主函数是vararg 所以是1
	IsVararg byte
	//寄存器数量  一个函数执行期间用多少个虚拟寄存器   Lua通常在编译期间统计好并且保存下来
	MaxStackSize byte
	//指令表  每条指令占4个字节  然后有几个指令
	Code []uint32
	//常量表 以1字节tag开头标记识别六种类型 nil,布尔值,整数,整数,浮点,字符串
	Constants []interface{}
	//Upvalues表  每个元素占2个字节
	Upvalues []Upvalue
	//函数原型表 长度为0
	Protos []*Prototype
	//行号表 每个指令对应的行号
	LineInfo []uint32
	//局部变量表 记录局部变量名  表中每个元素都包含变量名(按字符串储存) 和 起止指令索引(按cint储存)
	//如果主函数没有局部变量  则 长度是0
	LocVars []LocVar
	//Upvalue名列表 通常储存为_ENV
	UpvalueNames []string
	//如果编译时候加了-s
	//行号表 局部变量表 和 upvalue名列表  这三个储存的都是调试信息
	//Lua编译器就会在二进制chunk中把这三个表清空
}

//Upvalues表  每个元素占2个字节
type Upvalue struct {
	Instack byte //局部变量
	Idx     byte
}

//局部变量表 记录局部变量名  表中每个元素都包含变量名(按字符串储存) 和 起止指令索引(按cint储存)
type LocVar struct {
	VarName string
	StartPC uint32
	EndPC   uint32
}

//用于解析二进制chunk
func Undump(data []byte) *Prototype {
	reader := &reader{data}
	reader.checkHeader()        //跳过头部检验
	reader.readByte()           //跳过Upvalue数量
	return reader.readProto("") //读取函数原型
}

func IsBinaryChunk(data []byte) bool {
	return len(data) > 4 && string(data[:4]) == LUA_SIGNATURE
}
