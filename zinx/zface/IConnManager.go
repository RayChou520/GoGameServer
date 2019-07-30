package zface
/*链接管理抽象层*/
type IConnManager interface {
	AddConn(conn IConnection)
	RemoveConn(conn IConnection)
	GetConn(connId uint32)(IConnection,error)  //利用connID获取链接
	Len()  int   //获取当前链接
	ClearConn()  //删除并停止所有链接
}