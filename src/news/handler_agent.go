package news

import (
	"github.com/GaoShou012/frontier"
	"github.com/golang/glog"
	"google.golang.org/protobuf/proto"
	proto_im "im/proto/im"
	proto_news "im/proto/news"
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
func (a *agent) OnMessage(conn frontier.Conn, message *proto_im.Message) {
	msg := &proto_news.Subscribe{}
	err := proto.Unmarshal(message.Body, msg)
	if err != nil {
		glog.Errorln(err)
		return
	}
	a.subscribe(conn, msg)
}
func (a *agent) OnLeave(conn frontier.Conn) {
	index := conn.GetId() % a.parallel
	a.handlers[index].OnLeave(conn)
}

// 有新的消息，需要进行推送
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

func (a *agent) subscribe(conn frontier.Conn, subscribe *proto_news.Subscribe) {
	index := conn.GetId() % a.parallel
	a.handlers[index].OnSubscribe(conn, subscribe)
}
