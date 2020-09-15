package news

var Agent *agent

func init() {
	Agent = &agent{}
	Agent.OnInit()
}
