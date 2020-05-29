package binchunk

import "encoding/binary"
import "math"

//存放要被解析的二进制chunk数据
type reader struct {
	data []byte
}

//读取byte 1个字节 ,并且剩余前挪
func (self *reader) readByte() byte {
	b := self.data[0]
	self.data = self.data[1:]
	return b
}

//读取n个byte n个字节 ,并且剩余前挪
func (self *reader) readBytes(n uint) []byte {
	bytes := self.data[:n]
	self.data = self.data[n:]
	return bytes
}

//读取cint 4个字节 ,并且剩余前挪
func (self *reader) readUint32() uint32 {
	i := binary.LittleEndian.Uint32(self.data)
	self.data = self.data[4:]
	return i
}

//读取size_t 8个字节 ,并且剩余前挪
func (self *reader) readUint64() uint64 {
	i := binary.LittleEndian.Uint64(self.data)
	self.data = self.data[8:]
	return i
}

//读取Lua整数 8字节  映射会Go int64 ,并且剩余前挪
func (self *reader) readLuaInteger() int64 {
	return int64(self.readUint64())
}

//读取Lua浮点数 8字节  映射会Go float64 ,并且剩余前挪
func (self *reader) readLuaNumber() float64 {
	return math.Float64frombits(self.readUint64())
}

//读取string
func (self *reader) readString() string {
	size := uint(self.readByte()) //短字符串?
	if size == 0 {                //NUll字符串
		return ""
	}
	if size == 0xFF { //长字符串
		size = uint(self.readUint64())
	}
	bytes := self.readBytes(size - 1)
	return string(bytes)
}

func (self *reader) checkHeader() {
	if string(self.readBytes(4)) != LUA_SIGNATURE {
		panic("not a precompiled chunk!")
	}
	if self.readByte() != LUAC_VERSION {
		panic("version mismatch!")
	}
	if self.readByte() != LUAC_FORMAT {
		panic("format mismatch!")
	}
	if string(self.readBytes(6)) != LUAC_DATA {
		panic("corrupted!")
	}
	if self.readByte() != CINT_SIZE {
		panic("int size mismatch!")
	}
	if self.readByte() != CSIZET_SIZE {
		panic("size_t size mismatch!")
	}
	if self.readByte() != INSTRUCTION_SIZE {
		panic("instruction size mismatch!")
	}
	if self.readByte() != LUA_INTEGER_SIZE {
		panic("lua_Integer size mismatch!")
	}
	if self.readByte() != LUA_NUMBER_SIZE {
		panic("lua_Number size mismatch!")
	}
	if self.readLuaInteger() != LUAC_INT {
		panic("endianness mismatch!")
	}
	if self.readLuaNumber() != LUAC_NUM {
		panic("float format mismatch!")
	}
}

func (self *reader) readProto(parentSource string) *Prototype {
	source := self.readString()
	if source == "" {
		source = parentSource
	}
	return &Prototype{
		Source:          source,
		LineDefined:     self.readUint32(),
		LastLineDefined: self.readUint32(),
		NumParams:       self.readByte(),
		IsVararg:        self.readByte(),
		MaxStackSize:    self.readByte(),
		Code:            self.readCode(),
		Constants:       self.readConstants,
		Upvalues:        self.readUpValues(),
		Protos:          self.readProtos(source),
		LineInfo:        self.readLineInfo(),
		LocVars:         self.readLocVars(),
		UpValueNames:    self.readUpvalueName(),
	}
}

//获取指令表
func (self *reader) readCode() []uint32 {
	code := make([]uint32, self.readUint32())
	for i := range code {
		code[i] = self.readUint32()
	}
	return code
}

func (self *reader) readConstants() []interface{} {
	//TODO:
}
