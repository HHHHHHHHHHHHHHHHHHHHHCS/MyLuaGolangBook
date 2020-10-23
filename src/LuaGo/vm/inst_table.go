package vm

import . "LuaGo/api"

//默认批次50
const LFIELDS_PER_FLUSH = 50

func newTable(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1

	vm.CreateTable(Fb2int(b), Fb2int(c))
	vm.Replace(a)
}

//获取常量index c 的值 放入栈顶
//获取table b  和 栈顶的当作table索引
//得到 放到栈位置 a
func getTable(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1
	vm.GetRK(c)
	vm.GetTable(b)
	vm.Replace(a)
}

// b 索引  c val 放入栈顶
// 查找table a
// 设置值
func setTable(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1

	vm.GetRK(b)
	vm.GetRK(c)
	vm.SetTable(a)
}

//把栈堆的一系列数据  放入 table
//a table 栈位置  b 长度  c起始位置
//c 最大 9比特   所以保存批次数 这样就表达最大扩容数为 50*512 = 25600
func setList(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1

	//如果 c 大于零 则表达的批次是 c-1
	//否则 c 存在下一个指令 里面
	if c > 0 {
		c = c - 1
	} else {
		c = Instruction(vm.Fetch()).Ax()
	}

	bIsZero := b == 0
	if bIsZero {
		b = int(vm.ToInteger(-1)) - a - 1
		vm.Pop(1)
	}

	vm.CheckStack(1)
	idx := int64(c * LFIELDS_PER_FLUSH)
	for j := 1; j <= b; j++ {
		idx++
		vm.PushValue(a + j)
		vm.SetI(a, idx)
	}

	if bIsZero {
		for j := vm.RegisterCount() + 1; j <= vm.GetTop(); j++ {
			idx++
			vm.PushValue(j)
			vm.SetI(a, idx)
		}
		vm.SetTop(vm.RegisterCount())
	}
}
