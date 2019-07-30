package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/zface"
)

var GlobalObject *GlobalObj

type GlobalObj struct {
	TcpServer zface.IServer   //当前Zinx的全局server对象
	Host string		//当前服务器的主机id
	TcpPort int		//当前服务器监听端口号
	Name string	   //当前服务器名称
	Version string  //当前zinx版本号

	MaxPacketSize uint32  //数据包的最大值
	MaxConn int 	//最大链接数
	WorkerPoolSize uint32  //业务工作池的数量
	MaxWorkerTaskLen uint32  //业务工作worker对应负责的任务队列最大任务
}

func init()  {//初始变量，设置一些未加载json的默认值
	GlobalObject = &GlobalObj{
		Name:"ZinxServerApp",
	 	Version: "v0.4",
	 	TcpPort: 8999,
	 	Host: "0.0.0.0",

	 	MaxConn: 1200,
	 	MaxPacketSize: 4096,
		WorkerPoolSize: 10,
		MaxWorkerTaskLen: 1024,
	}
	GlobalObject.Reload()
}



func (g *GlobalObj) Reload() {
	data,err :=ioutil.ReadFile("conf/zinx.json")
	if err != nil{
		panic(err)
	}

	err = json.Unmarshal(data,&GlobalObject)

	if err!= nil{
		panic(err)
	}

}
