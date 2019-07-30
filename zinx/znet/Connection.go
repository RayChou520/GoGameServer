package znet

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/zface"
)

type Connection struct {
	TcpServer zface.IServer //当前conn属于哪个server，在conn初始化的时候添加即可

	Conn *net.TCPConn  //当前的链接

	ConnID uint32  //链接的ID

	IsClosed bool  //当前链接的状态

	ExitBuffChan chan bool //告知该链接已经退出/停止 channel

	MsgHandler zface.IMsgHandler  //该链接的处理方法router

	MsgChan chan []byte  //无缓冲管道，用于读、写两个goroutine之间的消息通信

	MsgBuffChan chan []byte   //有关冲管道，用于读、写两个goroutine之间的消息通信


	// ================================
	//链接属性
	property map[string]interface{}
	//保护链接属性修改的锁
	propertyLock sync.RWMutex
	// ================================
}

func (c *Connection) StartReader(){
	fmt.Println("[Debug] Reader Goroutine is running")
	//(defer后面必须跟函数，这个函数会在defer所在的函数执行完的时候调用，多个defer存在的时候后面的会先执行)
	defer fmt.Println("[Debug]",c.Conn.RemoteAddr().String(),"read exit")
	defer c.Stop()

	for{
		//创建拆包解包的对象
		dp := NewDataPack()

		//读取客户端的Msg head
		headData := make([]byte,dp.GetHeadLen())
		_,err := io.ReadFull(c.GetTcpConnection(),headData)
		if err != nil {
			fmt.Println("[Failed] receive buf failed")
			c.ExitBuffChan <- true
			break
		}

		//拆包，将得到的msgid和datalen放在msg中
		msg,err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("[Failed] unpack failed")
			c.ExitBuffChan <- true
			break
		}

		//根据dataLen 读取data，放在msg.data中
		var data []byte
		if msg.GetDataLen()>0{
			data = make([]byte,msg.GetDataLen())
			_,err := io.ReadFull(c.GetTcpConnection(),data)
			if err != nil {
				fmt.Println("[Failed] read msg data failed")
				c.ExitBuffChan <- true
				continue
			}
		}

		msg.SetData(data)
		//将链接和data绑定到一个请求的request数据中
		req := Request{
			conn:c,
			msg:msg,
		}

		if utils.GlobalObject.WorkerPoolSize>0 {
			//如果工作池开启，通过工作池机制将req交给worker处理
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			//从绑定好的消息和对应的处理方法中执行对应的Handle方法
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

func (c *Connection) StartWriter() {  //用户将消息发送给客户端
	fmt.Println("[Debug] Writer Goroutine is running")

	defer fmt.Println(c.RemoterAddr().String(),"[Debug] conn writer exit!")

	for{
		select {
		case data:= <-c.MsgChan:   //无缓存的msgChan
			if _,err :=c.Conn.Write(data); err!=nil {
				fmt.Println("[Failed] send data error:",err,"con writer exit")
				return
			}
		case data,ok:= <-c.MsgBuffChan:
			if ok{
				if _,err :=c.Conn.Write(data); err!=nil {
					fmt.Println("[Failed] send buff data error:", err, "con writer exit")
					return
				}
			}else {
				fmt.Println("[Debug] msgBuffChan is Closed")
				break
			}
			
		case <- c.ExitBuffChan:
			return //true,说明链接已经关闭了
		}
	}
}


//创建初始化链接的函数
func NewConnection(server zface.IServer,conn *net.TCPConn,connId uint32,msgHandler zface.IMsgHandler) *Connection{
	c:=&Connection{
		TcpServer:server,
		Conn: conn,
		ConnID:connId,
		IsClosed:false,
		ExitBuffChan:make(chan bool,1),
		MsgHandler:msgHandler,
		MsgChan:make(chan []byte),
		MsgBuffChan:make(chan []byte,3),

		property: make(map[string]interface{}), //对链接属性map初始化
	}
	c.TcpServer.GetConnMgr().AddConn(c)
	return  c
}

func (c *Connection) Start(){
	//1 开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()
	//2 开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()

	//==================
	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.TcpServer.CallOnConnStart(c)
	//==================

	for{  //循环检查在read数据的过程中，是否向通道chan中存放了true，如果有，说明发生了error，退出该函数
		select {
			case <- c.ExitBuffChan: //在未接收到true前，阻塞在这
				return //跳出函数
		}
	}
}	//启动链接(当监听到客户端请求后，初始化链接NewConnection后，开始执行读程序)

func (c *Connection) Stop() {
	if c.IsClosed == true {
		return
	}
	c.IsClosed = true

	//==================
	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.TcpServer.CallOnConnStop(c)
	//==================

	if err := c.Conn.Close(); err != nil{  //关闭socket该链接
		fmt.Println("[Failed] close conn failed",err)
	}

	c.ExitBuffChan <- true   //通知 读数据的业务，链接已读取完了，可以不在阻塞了可以关闭了
	c.TcpServer.GetConnMgr().RemoveConn(c) //删除conn从connManager中

	close(c.ExitBuffChan)  //关闭该链接管道
	close(c.MsgBuffChan)

} //关闭连接

func (c *Connection) GetTcpConnection() *net.TCPConn {
	return  c.Conn
}  //获取链接

func (c *Connection) GetConnectionID() uint32{
	return c.ConnID
}	//获取链接ID


func (c *Connection) RemoterAddr() net.Addr{
	return  c.Conn.RemoteAddr()
}  //获取远程链接的信息

func (c *Connection) SendMsg(msgId uint32, data []byte) error{
	if c.IsClosed == true {
		return errors.New("[Failed] connection closed when send msg")
	}
	//将data封包，并发送
	dp := NewDataPack()
	msg,err := dp.Pack(NewMsgPackage(msgId,data))
	if err!=nil{
		fmt.Println("[Failed] Pack error	msg id =",msgId)
		return errors.New("[Failed]Pack error msg")
	}

	//写回客户端
	c.MsgChan <- msg   //发送给channel供writer读取
	return nil
}

func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error{
	if c.IsClosed == true{
		return errors.New("[Failed] connection closed when send buff msg")
	}

	dp := NewDataPack()
	msg,err := dp.Pack(NewMsgPackage(msgId,data))
	if err!=nil{
		fmt.Println("[Failed] Pack error	msg id =",msgId)
		return errors.New("[Failed]Pack error msg")
	}

	c.MsgBuffChan <- msg
	return nil
} //添加带缓冲发 送消息接口


//设置链接属性
func (c *Connection) SetProperty(key string, value interface{}){
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}
//获取链接属性
func (c *Connection) GetProperty(key string)(interface{}, error){
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if value,ok := c.property[key]; ok {
		return value,nil
	}else {
		return nil,errors.New("[Debug] No property found")
	}
}
//移除链接属性
func (c *Connection) RemoveProperty(key string){
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property,key)
}

