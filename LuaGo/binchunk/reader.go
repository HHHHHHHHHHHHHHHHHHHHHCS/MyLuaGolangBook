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

func (self *reader) CheckHeader() {
	if string(self.readBytes(4)) != LUA_SIGNATURE {
		panic("not a precompiled chunk!")
	} else if self.readByte() != LUAC_VERSION {
		panic("version mismatch!")
	}
	else if self.readByte()!=LUAC_FORMAT{
		panic("format mismatch!")
	}
	else if string(self.readBytes(6)) != LUAC_DATA{
		panic("corrupted!")
	}
}
