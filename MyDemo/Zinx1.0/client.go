package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

func main() {
	fmt.Println("[Start] client is starting !")
	time.Sleep(3 * time.Second)

	con, err := net.Dial("tcp4", "127.0.0.1:8999")

	if err != nil {
		println("[Failed] client connection fail !", err)
		return
	}

	for {
		//发封包的message的消息
		dp := znet.NewDataPack()
		msg,_ := dp.Pack(znet.NewMsgPackage(1,[]byte("Zinx 0.6 client test message")))
		_, err := con.Write(msg)

		if err != nil {
			fmt.Println("[Failed] write con error", err)
			return
		}

		headData := make([]byte,dp.GetHeadLen())
		_,err = io.ReadFull(con,headData)
		if err != nil {
			fmt.Println("[Failed] read head error", err)
			break
		}

		msgHead,err :=dp.UnPack(headData)
		if err != nil {
			fmt.Println("[Failed] unpack error", err)
			return
		}

		if msgHead.GetDataLen()>0{
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte,msg.GetDataLen())

			_,err = io.ReadFull(con,msg.Data)
			if err != nil{
				fmt.Println("[Failed] server unpack data err")
				return
			}

			fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
		}

		time.Sleep(1 * time.Second)
	}

}
