package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/zface"
)

type MsgHandler struct {
	Apis map[uint32] zface.IRouter   //存放每一个msgId对应的处理方法

	WorkerPoolSize uint32 //业务工作worker池的数量

	TaskQueue []chan zface.IRequest   //负worker负责取任务的消息队列
}

func NewMsgHandler() *MsgHandler{
	return &MsgHandler{
		Apis:make(map[uint32]zface.IRouter),
		WorkerPoolSize:utils.GlobalObject.WorkerPoolSize,
		TaskQueue:make([]chan zface.IRequest,utils.GlobalObject.WorkerPoolSize),
	}
}

func (m *MsgHandler) DoMsgHandler(request zface.IRequest){
	handler,ok := m.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("[Failed] API msGid =",request.GetMsgId(),"is not found!")
		return
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}  //马上以非阻塞的方式处理

func (m *MsgHandler) AddRouter(msgId uint32,router zface.IRouter){
	if _,ok := m.Apis[msgId]; ok{
		panic("repeated api,msGid = " + strconv.Itoa(int(msgId)))
	}
	//判断未绑定过关系后，绑定
	m.Apis[msgId] = router
	fmt.Println("[Debug] add api msGid =",msgId)
} //为消息添加具体的处理逻辑

func (m *MsgHandler) StartWorkerPool() {
	//遍历需要启动worker的数量，依此启动
	for i:= 0; i < int(m.WorkerPoolSize); i++ {
		//一个worker被启动
		//给当前worker对应的任务队列开辟空间
		m.TaskQueue[i] = make (chan zface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}  //启动worker工作池

func (m *MsgHandler) SendMsgToTaskQueue(request zface.IRequest) {
	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则
	//得到需要处理此条连接的workerID
	workerID := request.GetConnection().GetConnectionID() % m.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnectionID()," request msgID=", request.GetMsgId(), "to workerID=", workerID)
	//将请求消息发送给任务队列
	m.TaskQueue[workerID] <- request
}  //将消息交给taskQueue，由worker处理

func (m *MsgHandler) StartOneWorker(workerID int, taskQueue chan zface.IRequest) {
	fmt.Println("[Debug] Worker ID = ", workerID, " is started.")
	//不断的等待队列中的消息
	for {
		select {
		//有消息则取出队列的Request，并执行绑定的业务方法
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		}
	}
}