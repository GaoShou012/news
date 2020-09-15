package news

var HandlerAgent *handlerAgent

var Agent *agent

func init() {
	Agent = &agent{}
	Agent.OnInit()
}
