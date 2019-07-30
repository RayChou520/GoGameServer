package main

import (
	"fmt"
	"zinx/zface"
	"zinx/znet"
)

//自定义路由(也叫自定义业务，让路由去处理相关的request)
type PingRouter struct {
	znet.BaseRouter  // 一定要先继承BaseRouter
}

func (p *PingRouter) Handle(request zface.IRequest)  {
	fmt.Println("[Debug] call PingRouter Handle")
	err := request.GetConnection().SendMsg(0,[]byte ("before...ping...ping\n"))

	if err != nil{
		fmt.Println("[Failed] call back ping ping ping error")
	}
}

type HelloRouter struct {
	znet.BaseRouter  // 一定要先继承BaseRouter
}

func (p *HelloRouter) Handle(request zface.IRequest)  {
	fmt.Println("[Debug] call PingRouter Handle")
	err := request.GetConnection().SendMsg(1,[]byte ("Hello...world...me\n"))

	if err != nil{
		fmt.Println("[Failed] call back ping ping ping error")
	}
}

//创建conn连接的时候执行
func DoConnectionBegin(conn zface.IConnection) {
	fmt.Println("DoConnectionBegin is Called ... ")

	//=============设置两个链接属性，在连接创建之后===========
	fmt.Println("Set conn Name, Home done!")
	conn.SetProperty("Name", "zzr")
	conn.SetProperty("Home", "hi go world")
	//===================================================

	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		fmt.Println(err)
	}
}
//连接conn断开的时候执行
func DoConnectionLost(conn zface.IConnection) {
	fmt.Println("DoConnectionLost is Called ... ")
	//============在连接销毁之前，查询conn的Name，Home属性=====
	if name, err:= conn.GetProperty("Name"); err == nil {
		fmt.Println("Conn Property Name = ", name)
	}
	if home, err := conn.GetProperty("Home"); err == nil {
		fmt.Println("Conn Property Home = ", home)
	}
	//===================================================
}


func main() {
	s := znet.NewServer()
	//注册链接hook的回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//配置路由服务
	s.AddRouter(0,&PingRouter{})
	s.AddRouter(1,&HelloRouter{})
	s.Server()
}
