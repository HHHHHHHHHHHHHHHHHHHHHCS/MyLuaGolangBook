package vm

import (
	. "LuaGo/api"
	. "LuaGo/state"
)

//拿到全局的表
func getTabUp(i Instruction,vm LuaVM){
	a,_,c :=i.ABC()
	a+=1

	vm.PushGlobalTable()
	vm.GetRK(c)
	vm.GetTable(-2)
	vm.Replace(a)
	vm.Pop(1)
}

//得到upvalue 的value
func getUpval(i Instruction,vm LuaVM)  {
	//虚拟机的upvalue 从0开始  Lua栈伪索引 从1开始
	a,b,_ :=i.ABC()
	a+=1
	b+=1

	vm.Copy(LuaUpvalueIndex(b),a)
}

func setUpval(i Instruction,vm LuaVM) {
	a,b,_:=i.ABC()
	a+=1
	b+=1


	vm.Copy(a,LuaUpvalueIndex(b))
}