package znet

type Message struct {
	Id  uint32
	DataLen uint32
	Data []byte
}

//初始化一个message的一个消息包
func NewMsgPackage(id uint32,data []byte) *Message{
	return &Message{
		Id:id,
		DataLen:uint32(len(data)),
		Data:data,
	}
}

func (msg *Message) GetDataLen() uint32{
	return msg.DataLen
}   //获取消息的长度
func (msg *Message) GetMsgId() uint32{
	return msg.Id
}  //获取消息的id
func (msg *Message) GetData() []byte{
	return msg.Data
}  //获取消息的长度

func (msg *Message) SetMsgId(uint32){}

func (msg *Message) SetData([]byte){}
//设计消息内容
func (msg *Message) SetDataLen(uint32){}	//设置消息数据段的长度