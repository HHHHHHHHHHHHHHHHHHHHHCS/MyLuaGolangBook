package vm

import . "LuaGo/api"

//运算符相关的指令
func _binaryArith(i Instruction, vm LuaVM, op ArithOp) {
	a,b,c:=i.ABC()
	a+=1

	vm.GetRK(b)
	vm.GetRK(c)
	vm.Arith(op)
	vm.Replace(a)
}
