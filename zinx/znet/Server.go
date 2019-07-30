package znet

import (
	"fmt"
	"net"
	"time"
	"zinx/utils"
	"zinx/zface"
)

type Server struct {
	Name      string
	IpVersion string
	Ip        string
	Port      int
	MaxConn	  int
	msgHandler    zface.IMsgHandler  //当前Server的消息管理模块，用来绑定MsgId和对应的处理方法
	ConnMgr  zface.IConnManager //当前server的链接管理器

	OnConnStart func(conn zface.IConnection)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn zface.IConnection)
}

func (s *Server) Start() {
	go func() {
		fmt.Printf("[Start] Server Listening of Ip: %s:%d, MaxConn: %d,is starting...\n", s.Ip, s.Port,s.MaxConn)
		fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPacketSize:%d\n",
			utils.GlobalObject.Version,
			utils.GlobalObject.MaxConn,
			utils.GlobalObject.MaxPacketSize)

		s.msgHandler.StartWorkerPool()  //启动worker工作池机制
		//解析一个tcp的addr
		addr, err := net.ResolveTCPAddr(s.IpVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))

		if err != nil {
			fmt.Println("[Error] Resolve tcp addr error ", err)
			return
		}

		//监听服务器的地址
		listenner, err := net.ListenTCP(s.IpVersion, addr)

		if err != nil {
			fmt.Println("[Error] Listening tcp addr error ", err)
			return
		}

		fmt.Println("[Success] Start server", s.Name, "success now is listening...")

		//阻塞的等待客户端的链接，处理客户端链接业务
		var cid uint32
		cid = 0
		for {// 启动sever链接的服务
			//如果有客户端链接过来，建立连接，阻塞返回
			con, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("[Error] Accept error ", err)
				continue
			}

			//如果超过最大链接，那么关闭该链接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				_ =con.Close()
				continue
			}


			dealConn :=NewConnection(s,con,cid,s.msgHandler)
			cid ++

			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name " , s.Name)
	s.ConnMgr.ClearConn()
}

func (s *Server) Server() {
	s.Start()

	//TODO

	for{
		time.Sleep(10*time.Second)
	}
}

func (s *Server) AddRouter(msgId uint32 ,router zface.IRouter){
	s.msgHandler.AddRouter(msgId,router)
	fmt.Println("[Debug] Add router success")
}


func (s *Server) GetConnMgr() zface.IConnManager{   //得到链接管理
	return s.ConnMgr
}


func NewServer() zface.IServer {
	utils.GlobalObject.Reload() // 加载服务器配置文件

	s := &Server{
		Name:      utils.GlobalObject.Name,
		IpVersion: "tcp4",
		Ip:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		MaxConn:   utils.GlobalObject.MaxConn,
		msgHandler:    NewMsgHandler(),
		ConnMgr: NewConnManager(),
	}
	return s
}

//设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(HookFunc func (zface.IConnection)){
	s.OnConnStart = HookFunc
}
//设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(HookFunc func (zface.IConnection)){
	s.OnConnStop = HookFunc
}
//调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn zface.IConnection){
	if s.OnConnStart != nil {
		fmt.Println("--->CallOnConnStart....")
		s.OnConnStart(conn)   //调用的是server中的OnConnStart属性（这个属性的类型是一个函数）,调用即执行
	}
}
//调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn zface.IConnection){
	if s.OnConnStop != nil {
		fmt.Println("--->CallOnConnStop....")
		s.OnConnStop(conn)   //调用的是server中的OnConnStart属性（这个属性的类型是一个函数）,调用即执行
	}
}