package news

import (
	"github.com/golang/glog"
	"google.golang.org/protobuf/proto"
	proto_news "im/proto/news"
	"im/src/frontier"
	"runtime"
)

type agent struct {
	parallel int
	handlers []*handler
}

func (a *agent) OnInit() {
	a.parallel = runtime.NumCPU()
	a.handlers = make([]*handler, a.parallel)
	for i := 0; i < a.parallel; i++ {
		a.handlers[i].OnInit()
	}
}
func (a *agent) OnSubscribe(conn frontier.Conn, subscribe *proto_news.Subscribe) {
	index := conn.GetId() % a.parallel
	a.handlers[index].OnSubscribe(conn, subscribe)
}
func (a *agent) OnLeave(conn frontier.Conn) {
	index := conn.GetId() % a.parallel
	a.handlers[index].OnLeave(conn)
}
func (a *agent) OnNews(data []byte) {
	item := &proto_news.NewsItem{}
	err := proto.Unmarshal(data, item)
	if err != nil {
		glog.Errorln(err)
		return
	}
	for i := 0; i < a.parallel; i++ {
		a.handlers[i].OnNews(item)
	}
}
