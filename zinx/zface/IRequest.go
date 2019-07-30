package zface

type IRequest interface {
	GetConnection() IConnection 	//获取请求链接的信息
	GetData() []byte   //获取请求数据的消息
	GetMsgId() uint32
}