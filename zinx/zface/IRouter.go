package zface

type IRouter interface {
	PreHandle (request IRequest) //处理业务之前的构造方法
	Handle (request IRequest)
	PostHandle(request IRequest) //处理业务之后的构造方法
}