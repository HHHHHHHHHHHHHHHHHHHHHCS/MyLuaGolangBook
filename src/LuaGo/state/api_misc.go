package state

//暂时只考虑 字符串的长度
func (self *luaState) Len(idx int) {
	val := self.stack.get(idx)
	if s, ok := val.(string); ok {
		self.stack.push(int64(len(s)))
	} else {
		panic("length error!")
	}
}

//栈顶弹出N个值 然后for循环 转换为 string 进行拼接 然后推入栈顶
//如果栈长度为0 则推入""
func (self *luaState) Concat(n int) {
	if n == 0 {
		self.stack.push("")
	} else if n >= 2 {
		for i := 1; i < n; i++ {
			if self.IsString(-1) && self.IsString(-2) {
				s2 := self.ToString(-1)
				s1 := self.ToString(-2)
				self.stack.pop()
				self.stack.pop()
				self.stack.push(s1 + s2)
				continue
			}
			panic("concatenation error!")
		}
	}
	//if n == 1 , do nothing
}
