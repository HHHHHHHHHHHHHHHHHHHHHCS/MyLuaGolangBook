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

//把操作数A索引换成操作数B bool值  如果布尔值!=0 则跳过下一个指令
func loadBool(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	//推入栈顶 bool 值   在移除栈顶 换到位置a
	vm.PushBoolean(b != 0)
	vm.Replace(a)
	if c != 0 {
		vm.AddPC(1)
	}

}

//读取一般索引的常量 压入栈位置a
func loadK(i Instruction, vm LuaVM) {
	//bx 是 32-14=18bit 最大无符号整值 262143
	//常量表长度通常不会超过这个大小
	a, bx := i.ABx()
	a += 1

	//放入常量表bx到栈顶    再移除栈顶 放到索引a
	vm.GetConst(bx)
	vm.Replace(a)
}

//读取很大索引的常量  因为Lua是数据描述语言 有时候会出现这种情况
func loadKx(i Instruction, vm LuaVM) {
	//用两个指令  第一个指令是索引  第二个指令是常量索引
	a, _ := i.ABx()
	a += 1
	//ax 这样就占 32-6=26bit 最大无符号整数67108864
	ax := Instruction(vm.Fetch()).Ax()

	vm.GetConst(ax)
	vm.Replace(a)
}
