package zface

import "net"

type IConnection interface {
	Start()	//启动链接

	Stop()  //关闭连接

	GetTcpConnection() *net.TCPConn   //获取链接

	GetConnectionID() uint32	//获取链接ID

	RemoterAddr() net.Addr  //获取远程链接的信息

	SendMsg(msgId uint32, data []byte) error

	SendBuffMsg(msgId uint32, data []byte) error //添加带缓冲发 送消息接口


	//设置链接属性
	SetProperty(key string, value interface{})
	//获取链接属性
	GetProperty(key string)(interface{}, error)
	//移除链接属性
	RemoveProperty(key string)
}
