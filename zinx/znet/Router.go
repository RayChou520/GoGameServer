package znet

import "zinx/zface"

type BaseRouter struct {

}


func (b *BaseRouter) PreHandle (request zface.IRequest){

} //处理业务之前的构造方法
func (b *BaseRouter) Handle (request zface.IRequest){

}
func (b *BaseRouter) PostHandle(request zface.IRequest){

} //处理业务之后的构造方法


