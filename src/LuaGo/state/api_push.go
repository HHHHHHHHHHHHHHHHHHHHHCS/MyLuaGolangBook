package state

import . "LuaGo/api"

//lua值由外部推入栈顶
func (self *luaState) PushNil()                    { self.stack.push(nil) }
func (self *luaState) PushBoolean(b bool)          { self.stack.push(b) }
func (self *luaState) PushInteger(n int64)         { self.stack.push(n) }
func (self *luaState) PushNumber(n float64)        { self.stack.push(n) }
func (self *luaState) PushString(s string)         { self.stack.push(s) }
func (self *luaState) PushGoFunction(f GoFunction) { self.stack.push(newGoClosure(f)) }
