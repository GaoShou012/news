package news

import (
	"github.com/GaoShou012/frontier"
	"github.com/golang/glog"
	"google.golang.org/protobuf/proto"
	proto_im "im/proto/im"
	proto_news "im/proto/news"
	"runtime"
	"sync"
)

type agent struct {
	parallel int
	handlers []*handler

	connections     map[int]frontier.Conn
	subscribeOnConn map[int][]string

	mutex       sync.Mutex
	clients     map[int]*Client
	onSubscribe chan []string
	onCancel    chan []string
}

func (a *agent) OnInit() {
	a.parallel = runtime.NumCPU()
	a.handlers = make([]*handler, a.parallel)
	for i := 0; i < a.parallel; i++ {
		h := &handler{}
		h.OnInit()
		a.handlers[i] = h
	}

	for {
		select {
		case subList := <-a.onSubscribe:

			break
		case subList := <-a.onCancel:
			break
		}
	}
}
func (a *agent) OnMessage(conn frontier.Conn, message *proto_im.Message) {
	// 反序列化消息
	msg := &proto_news.Subscribe{}
	err := proto.Unmarshal(message.Body, msg)
	if err != nil {
		glog.Errorln(err)
		return
	}

	// 订阅频道数量限制
	if len(msg.Channels) <= 0 || len(msg.Channels) >= 255 {
		return
	}

	// 去重，如果出现重复，不进行操作
	temp := make(map[string]bool)
	for _, channelName := range msg.Channels {
		_, ok := temp[channelName]
		if ok {
			break
		}
	}

	a.mutex.Lock()
	defer a.mutex.Unlock()

	// 如果客户端不存在，创建客户端对象，订阅新的频道
	// 如果客户端已经存在，取消先前的订阅，然后，订阅新的频道
	clientId := conn.GetId()
	cli, ok := a.clients[clientId]
	if ok {
		a.onCancel <- cli.Subscribe
	} else {
		cli = &Client{
			Conn:      conn,
			News:      make(map[string]*proto_news.NewsItem),
			Subscribe: nil,
		}
		a.clients[clientId] = cli
	}

	cli.Subscribe = msg.Channels
	a.onSubscribe <- msg.Channels

	a.subscribe(conn, msg)
}
func (a *agent) OnLeave(conn frontier.Conn) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	clientId := conn.GetId()
	cli, ok := a.clients[clientId]
	if !ok {
		return
	}
	a.onCancel <- cli.Subscribe
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
