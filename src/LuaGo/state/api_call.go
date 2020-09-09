package state

import (
	"LuaGo/binchunk"
	"fmt"
)

//mode "b" 二进制 "t" 文本 "bt" 二进制或者文本
//0 表示加载成功   非0 不成功
func (self *luaState) Load(chunk []byte, chunkName, mode string) int {
	proto := binchunk.Undump(chunk)
	c := newLuaClosure(proto)
	self.stack.push(c)
	return 0
}

//正常闭包在栈顶-nArgs的位置  所以 nArgs 是参数位置 同时也暗示告诉包的位置
//nResults 返回多少个参数
func (self *luaState) Call(nArgs, nResults int) {
	val := self.stack.get(-(nArgs + 1))
	if c, ok := val.(*closure); ok {
		fmt.Printf("call %s<%d,%d>\n", c.proto.Source,
			c.proto.LineDefined, c.proto.LastLineDefined)
		self.callLuaClosure(nArgs, nResults, c)
	} else{
		panic("not function!")
	}
}
