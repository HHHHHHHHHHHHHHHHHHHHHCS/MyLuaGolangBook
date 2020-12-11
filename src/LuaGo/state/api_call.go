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
	//设置_ENV
	if len(proto.Upvalues) > 0 {
		env := self.registry.get(LUA_RIDX_GLOBALS)
		c.upvals[0] = &upvalue{&env}
	}
	return 0
}

//正常闭包在栈顶-nArgs的位置  所以 nArgs 是参数位置 同时也暗示告诉包的位置
//nResults 返回多少个参数  如果是-1  则返回全部参数
//如果栈顶不是方法 则看看元表 __call
func (self *luaState) Call(nArgs, nResults int) {
	val := self.stack.get(-(nArgs + 1))

	c, ok := val.(*closure)
	//如果有有元表方法__call 则覆盖c和ok 把当前表也当作数据传入参数
	if !ok {
		if mf := getMetafield(val, "__call", self); mf != nil {
			if c, ok = mf.(*closure); ok {
				self.stack.push(val)
				self.Insert(-(nArgs + 2))
				nArgs += 1
			}
		}
	}
	if ok {
		//如果有proto 则是lua方法 否则是go方法
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

	newStack := newLuaStack(nRegs+LUA_MINSTACK, self)
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
	newStack := newLuaStack(nArgs+LUA_MINSTACK, self)
	newStack.closure = c

	if nArgs > 0 {
		args := self.stack.popN(nArgs)
		newStack.pushN(args, nArgs)
	}
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

func (self *luaState) PCall(nArgs, nResults, msgh int) (status int) {
	caller := self.stack
	status = LUA_ERRRUN
	// 异常捕获机制
	defer func() {
		//recover是一个从panic恢复的内建函数。
		//recover只有在defer的函数里面才能发挥真正的作用
		if err := recover(); err != nil {
			for self.stack != caller {
				self.popLuaStack()
			}
			self.stack.push(err)
		}
	}()

	self.Call(nArgs, nResults)
	//如果call 失败了  则直接 跳到  return 在 defer   最后输出
	status = LUA_OK
	return
}
