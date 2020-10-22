package vm

import (
	. "LuaGo/api"
)

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

//如果要找的upvalue 是table
func getTabUp(i Instruction,vm LuaVM) {
	a,b,c := i.ABC()
	a+=1
	b+=1
	vm.GetRK(c)
	vm.GetTable(LuaUpvalueIndex(b))
	vm.Replace(a)
}

//给upvalue 的table 设置值
func setTabUp(i Instruction,vm LuaVM) {
	a,b,c := i.ABC()
	a+=1

	vm.GetRK(b)
	vm.GetRK(c)
	vm.SetTable(LuaUpvalueIndex(a))
}