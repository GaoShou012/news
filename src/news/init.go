package news

//var Agent *agent

var DefaultCodec *Codec

func init() {
	//Agent = &agent{}
	//Agent.OnInit()

	DefaultCodec = &Codec{}
	DefaultCodec.Run()
}
