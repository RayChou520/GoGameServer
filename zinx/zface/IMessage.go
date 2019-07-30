package zface

//消息封装到message中，实现抽象层的接口

type IMessage interface {
	GetDataLen() uint32   //获取消息的长度
	GetMsgId() uint32  //获取消息的id
	GetData() []byte  //获取消息的长度

	SetMsgId(uint32)
	SetData([]byte)		//设计消息内容
	SetDataLen(uint32)	//设置消息数据段的长度
}

