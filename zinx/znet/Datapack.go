package znet

import (
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
	"zinx/utils"
	"zinx/zface"
)

type DataPack struct {}

//封包拆包实例初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32{
	//Id uint32(4字节) + DataLen uint32(4字节)
	return 8
} //获取包头长度的方法

func (d *DataPack) Pack(msg zface.IMessage)([]byte,error) {
	//创建一个存放bytes字节的缓冲
	databuff := bytes.NewBuffer([]byte{})
	//写datalen
	if err := binary.Write(databuff,binary.LittleEndian,msg.GetDataLen()); err!=nil{
		return nil,err
	}
	//写msgid
	if err := binary.Write(databuff,binary.LittleEndian,msg.GetMsgId()); err!=nil{
		return nil,err
	}
	//写data数据
	if err := binary.Write(databuff,binary.LittleEndian,msg.GetData()); err!=nil{
		return nil,err
	}
		return databuff.Bytes(),nil
} //封包方法

func (d *DataPack) UnPack(binaryData []byte)(zface.IMessage,error) {
	//创建一个输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head信息，得到datalen和msgid
	msg := &Message{}

	if err := binary.Read(dataBuff,binary.LittleEndian,&msg.DataLen); err!=nil{
		return nil,err
	}

	if err := binary.Read(dataBuff,binary.LittleEndian,&msg.Id); err!=nil{
		return nil,err
	}

	//判断datalen的长度是否超出我们允许的最大包长度
	if utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > utils.GlobalObject.MaxPacketSize {
		return nil,errors.New("Too large msg data received")
	}

	//以上是把head的数据包拆出来就可以了，然后再通过head的长度，在从conn中取一次数据

	return msg,nil

} //拆包方法

