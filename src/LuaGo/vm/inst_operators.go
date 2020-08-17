package vm

import . "LuaGo/api"

//运算符相关的指令

//二元运算
func _binaryArith(i Instruction, vm LuaVM, op ArithOp) {
	a, b, c := i.ABC()
	a += 1

	vm.GetRK(b)
	vm.GetRK(c)
	vm.Arith(op)
	vm.Replace(a)
}

//一元自运算
func _unaryArith(i Instruction, vm LuaVM, op ArithOp) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	vm.PushValue(b)
	vm.Arith(op)
	vm.Replace(a)
}

//+
func add(i Instruction, vm LuaVM) {
	_binaryArith(i, vm, LUA_OPADD)
}

//-
func sub(i Instruction, vm LuaVM) {
	_binaryArith(i, vm, LUA_OPSUB)
}

//*
func mul(i Instruction, vm LuaVM) {
	_binaryArith(i, vm, LUA_OPMUL)
}

//%
func mod(i Instruction, vm LuaVM) {
	_binaryArith(i, vm, LUA_OPMOD)
}

//^
func pow(i Instruction, vm LuaVM) {
	_binaryArith(i, vm, LUA_OPPOW)
}

// /
func div(i Instruction, vm LuaVM) {
	_binaryArith(i, vm, LUA_OPDIV)
}

// // 整除
func idiv(i Instruction, vm LuaVM) {
	_binaryArith(i, vm, LUA_OPIDIV)
}

// & bool 且
func band(i Instruction, vm LuaVM) {
	_binaryArith(i, vm, LUA_OPBAND)
}

// | bool 或
func bor(i Instruction, vm LuaVM) {
	_binaryArith(i, vm, LUA_OPBOR)
}

// ~ bool 异或
func bxor(i Instruction, vm LuaVM) {
	_binaryArith(i, vm, LUA_OPBXOR)
}

// << 整数 左位移
func shl(i Instruction, vm LuaVM) {
	_binaryArith(i, vm, LUA_OPSHL)
}

// >> 整数 右位移
func shr(i Instruction, vm LuaVM) {
	_binaryArith(i, vm, LUA_OPSHR)
}

// - 自负数
func unm(i Instruction, vm LuaVM) {
	_unaryArith(i, vm, LUA_OPUNM)
}

// ~ 自取反
func bnot(i Instruction, vm LuaVM) {
	_unaryArith(i, vm, LUA_OPBNOT)
}

//得到某个栈的值的长度(如string) 在塞入某个位置
func length(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	vm.Len(b)
	vm.Replace(a)
}

//for 循环字符串拼接  把结果替换到 某个位置
//字符串拼接需要检测长度  不然容易栈溢出
func concat(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1
	c += 1

	n := c - b + 1
	vm.CheckStack(n)
	for i := b; i <= c; i++ {
		vm.PushValue(i)
	}

	vm.Concat(n)
	vm.Replace(a)
}

//比较 如果失败 则跳过下一条指令
func _compare(i Instruction, vm LuaVM, op CompareOp) {
	a, b, c := i.ABC()

	vm.GetRK(b)
	vm.GetRK(c)
	//a!=0 true   result!=true -> false ->jump 1 pc
	if vm.Compare(-2, -1, op) != (a != 0) {
		vm.AddPC(1)
	}
	vm.Pop(2)
}

//==
func eq(i Instruction, vm LuaVM) {
	_compare(i, vm, LUA_OPEQ)
}

//<
func lt(i Instruction, vm LuaVM) {
	_compare(i, vm, LUA_OPLT)
}

//<=
func le(i Instruction, vm LuaVM) {
	_compare(i, vm, LUA_OPLE)
}

//读取栈索引b bool 取反 插入 栈索引a
func not(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	vm.PushBoolean(!vm.ToBoolean(b))
	vm.Replace(a)
}

//索引b的bool 值 是否和c 一样
//如果一样 则把 b 的值 赋值到栈索引a
//否则跳过下一条指令
func testSet(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	if vm.ToBoolean(b) == (c != 0) {
		vm.Copy(b, a)
	} else {
		vm.AddPC(1)
	}
}

//读取栈索引a  和 c 布尔比较
//如果失败 则跳过下一条指令
func test(i Instruction, vm LuaVM) {
	a, _, c := i.ABC()
	a += 1

	if vm.ToBoolean(a) != (c != 0) {
		vm.AddPC(1)
	}
}

