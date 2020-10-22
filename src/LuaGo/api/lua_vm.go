package api

type LuaVM interface {
	LuaState             //golang 的成员
	PC() int             //返回当前PC(仅测试用) PC program counter 程序计数器
	AddPC(n int)         //修改PC(用于实现跳转指令)
	Fetch() uint32       //取出当前的指令 , 将PC指向下一条指令
	GetConst(idx int)    //将指定常量推入栈顶
	GetRK(rk int)        //将指定常量或栈值推入栈顶
	RegisterCount() int  //读取常量的数量
	LoadVararg(n int)    //读取可变参数
	LoadProto(idx int)   //读取索引 方法
	CloseUpvalues(a int) //把upvalue 进行闭包
}
