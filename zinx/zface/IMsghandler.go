package zface

/*消息管理抽象层*/
type IMsgHandler interface {
	DoMsgHandler(request IRequest)  //马上以非阻塞的方式处理
	AddRouter(msgId uint32,router IRouter)  //为消息添加具体的处理逻辑

	StartWorkerPool()   //启动worker工作池
	SendMsgToTaskQueue(request IRequest)   //将消息交给taskQueue，由worker处理
}
