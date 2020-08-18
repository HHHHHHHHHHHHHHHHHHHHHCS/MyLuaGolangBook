package vm

import . "LuaGo/api"

//其他指令
//如:MOVE JMP

//移动指令 把原索引的变量拷贝到一个新索引位置  原来的不变
//2^8=256 lua的局部变量不能超过255个 否则会超出索引不能编译
func move(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	//go 索引是0开始  lua索引栈是1开始
	a += 1
	b += 1
	vm.Copy(b, a)
}

//跳转指令
func jmp(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	vm.AddPC(sBx)
	if a != 0 {
		panic("TODO: jmp!")
	}
}
