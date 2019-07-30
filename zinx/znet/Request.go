package znet

import "zinx/zface"

type Request struct {
	conn zface.IConnection //已经和客户端建立好的连接
	msg zface.IMessage  //客户端请求的数据
}

func (r *Request) GetConnection() zface.IConnection{
	return r.conn
} 	//获取请求链接的信息

func (r *Request) GetData() []byte{
	return r.msg.GetData()
}   //获取请求数据的消息

func (r *Request) GetMsgId() uint32{
	return r.msg.GetMsgId()
}