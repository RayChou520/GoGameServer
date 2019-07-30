package zface

type IServer interface { //定义一个服务器接口

	Start()  //启动服务器方法

	Stop()  //关闭服务器方法

	Server() //开启服务器业务方法

	AddRouter(msgId uint32 ,router IRouter) //路由功能:给当前服务注册一个路由业务的方法，供客户端链接处理使用

	GetConnMgr() IConnManager   //得到链接管理


	//设置该Server的连接创建时Hook函数
	SetOnConnStart(func (IConnection))
	//设置该Server的连接断开时的Hook函数
	SetOnConnStop(func (IConnection))
	//调用连接OnConnStart Hook函数
	CallOnConnStart(conn IConnection)
	//调用连接OnConnStop Hook函数
	CallOnConnStop(conn IConnection)
}
