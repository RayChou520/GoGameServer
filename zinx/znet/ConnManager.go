package znet

import (
	"fmt"
	"github.com/pkg/errors"
	"sync"
	"zinx/zface"
)

type ConnManager struct {
	Connections map[uint32]zface.IConnection  //管理的连接信息
	connLock sync.RWMutex  //读写连接的的读写锁
}

//创建一个连接管理
func NewConnManager() *ConnManager{
	return &ConnManager{
		Connections:make(map[uint32] zface.IConnection),
	}
}

func (cm *ConnManager) AddConn(conn zface.IConnection){
	//保存共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	//将conn链接添加到connManager中
	cm.Connections[conn.GetConnectionID()] = conn
	fmt.Println("[Debug] connection add to ConnManager successfully: conn num = ", cm.Len())
}

func (cm *ConnManager) RemoveConn(conn zface.IConnection){
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	delete(cm.Connections,conn.GetConnectionID())
	fmt.Println("connection Remove ConnID=",conn.GetConnectionID(), "successfully: conn num = ", cm.Len())
}

func (cm *ConnManager) GetConn(connId uint32)(zface.IConnection,error){
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	if conn,ok :=cm.Connections[connId]; ok{
		return conn,nil
	}else {
		return nil,errors.New("[Debug] connection not found")
	}

}//利用connID获取链接

func (cm *ConnManager) Len() int {
	return len(cm.Connections)
}  //获取当前链接数量

func (cm *ConnManager) ClearConn(){
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	//停止并删除全部的连接信息
	for connID,conn := range cm.Connections {
		conn.Stop()
		delete(cm.Connections,connID)
	}
	fmt.Println("[Debug] Clear All Connections successfully: conn num = ", cm.Len())
}  //删除并停止所有链接