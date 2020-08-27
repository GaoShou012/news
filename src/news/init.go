package news

var HandlerAgent *handlerAgent

func init() {
	HandlerAgent = &handlerAgent{}
	HandlerAgent.init()
}
