package LuaGo

type header struct {
	//独特的签名 用于识别文件 ESCLua 的十六进制 0x1B4C7564
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
	luacInt byte
	//lua浮点数 n(8)个字节存放 370.5 目的检测二进制chunk的浮点数格式 如IEEE754
	luacNum float64
}

type BinaryChunk struct {
	header                  //头部
	sizeUpvalues byte       //主函数upvalue的数量
	mainFunc     *Prototype //主函数原型
}
