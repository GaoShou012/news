package room

var (
	Agent Rooms
	Codec codec
)

func init() {
	Codec.init(10000)

	Agent.Init()
	Sender.Init()
	SyncRecord.Init(100)
}
