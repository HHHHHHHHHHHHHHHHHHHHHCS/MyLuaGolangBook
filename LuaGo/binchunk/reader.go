package binchunk

import "encoding/binary"
import "math"

//存放要被解析的二进制chunk数据
type reader struct {
	data []byte
}

//读取第一个字节 ,并且剩余前挪
func (self *reader) readByte() byte {
	b := self.data[0]
	self.data = self.data[1:]
	return b
}
