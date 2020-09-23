package state

import (
	. "LuaGo/api"
	"LuaGo/binchunk"
	"LuaGo/vm"
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
		//如果有proto 则是lua方法 否则是gofangfa
		if c.proto != nil {
			self.callLuaClosure(nArgs, nResults, c)
		} else {
			self.callGoClosure(nArgs, nResults, c)
		}
	} else {
		panic("not function!")
	}
}

func (self *luaState) callLuaClosure(nArgs, nResults int, c *closure) {
	nRegs := int(c.proto.MaxStackSize)
	nParams := int(c.proto.NumParams)
	isVararg := c.proto.IsVararg == 1

	newStack := newLuaStack(nRegs + LUA_MINSTACK)
	newStack.closure = c

	//把参数写入新的栈   并且 写入可变参数
	funcAndArgs := self.stack.popN(nArgs + 1)
	newStack.pushN(funcAndArgs[1:], nParams)
	newStack.top = nRegs
	if nArgs > nParams && isVararg {
		newStack.varargs = funcAndArgs[nParams+1:]
	}

	//把栈写入当前的线程  并且执行
	self.pushLuaStack(newStack)
	self.runLuaClosure()
	self.popLuaStack()

	//把多返回值 写入当前的栈
	if nResults != 0 {
		results := newStack.popN(newStack.top - nRegs)
		self.stack.check(len(results))
		self.stack.pushN(results, nResults)
	}
}

func (self *luaState) runLuaClosure() {
	for {
		inst := vm.Instruction(self.Fetch())
		inst.Execute(self)
		if inst.Opcode() == vm.OP_RETURN {
			break
		}
	}
}

func (self *luaState) callGoClosure(nArgs, nResults int, c *closure) {
	newStack := newLuaStack(nArgs + LUA_MINSTACK)
	newStack.closure = c

	args := self.stack.popN(nArgs)
	newStack.pushN(args, nArgs)
	self.stack.pop()

	self.pushLuaStack(newStack)
	r := c.goFunc(self)
	self.popLuaStack()

	if nResults != 0 {
		results := newStack.popN(r)
		self.stack.check(len(results))
		self.stack.pushN(results, nResults)
	}
}
