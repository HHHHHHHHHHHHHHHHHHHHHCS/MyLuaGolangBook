package vm

import . "LuaGo/api"

//加载相关的指令

//让栈内连续n个为nil  初始化用
//因为lua编译的时候 寄存器数量预先算好了
//然后执行阶段 SetTop() 保留好了栈空位  现在初始化
func loadNil(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1 //goApi 索引0   lua栈索引1

	//从顶部压入nil copy n个  再弹出栈顶
	vm.PushNil()
	for i := a; i <= a+b; i++ {
		vm.Copy(-1, i)
	}
	vm.Pop(1)
}
