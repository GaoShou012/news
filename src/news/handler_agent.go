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

	channels map[string]int

	mutex          sync.Mutex
	connections    map[int]frontier.Conn
	connectionsSub map[int][]string
	onSubscribe    chan []string
	onCancelSub    chan []string
	onPulling      chan []string
}

func (a *agent) OnInit() {
	a.parallel = runtime.NumCPU()
	a.handlers = make([]*handler, a.parallel)
	for i := 0; i < a.parallel; i++ {
		h := &handler{}
		h.OnInit()
		a.handlers[i] = h
	}

	a.channels = make(map[string]int)
	a.connections = make(map[int]frontier.Conn)
	a.connectionsSub = make(map[int][]string)
	a.onSubscribe = make(chan []string)
	a.onCancelSub = make(chan []string)

	go func() {
		// 理论上任务顺序是，先订阅，后取消。
		// 实际上有可能任务执行的时候，订阅事件还没有被消费，取消事件就已经优先消费
		// 所以在取消事件中 count-- 如果是 -1 不会产生取消任务，如果是 0 就产生
		// 如果取消事件中count-- equal -1 订阅事件在还排队中，所以，在订阅处理中，count++ 的时候，需要 1 才产生订阅任务
		// 最终结果 count-- equal -1 不产生任务，count++ equal 0 不产生任务
		// 必须是： count-- equal 0 生成任务，count++ equal 1 产生任务
		for {
			select {
			case subList := <-a.onSubscribe:
				for _, channelName := range subList {
					count, ok := a.channels[channelName]
					if !ok {
						count = 0
					}
					count++
					a.channels[channelName] = count

					if count == 1 {
						// 订阅
					}
				}
				break
			case subList := <-a.onCancelSub:
				for _, channelName := range subList {
					count, ok := a.channels[channelName]
					if !ok {
						count = 0
					}
					count--
					a.channels[channelName] = count

					if count == 0 {
						// 取消订阅
					}
				}
				break
			}
		}
	}()

	for i := 0; i < runtime.NumCPU()*10; i++ {
		go func() {
			for {
				task := <-a.onPulling

			}
		}()
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
			return
		}
	}

	a.mutex.Lock()
	defer a.mutex.Unlock()

	// 如果终端连接不存在，保存连接，订阅新的频道
	// 如果终端连接经存在，取消先前的订阅，然后，订阅新的频道
	connId := conn.GetId()
	_, ok := a.connections[connId]
	if !ok {
		a.connections[connId] = conn
	} else {
		a.onCancelSub <- a.connectionsSub[connId]
	}
	// Save the subscribe
	a.connectionsSub[connId] = msg.Channels
	// Subscribe the news
	a.onSubscribe <- a.connectionsSub[connId]
	// Subscribe the news in the handler
	a.handler(conn).OnSubscribe(conn, msg)
}
func (a *agent) OnLeave(conn frontier.Conn) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	connId := conn.GetId()
	_, ok := a.connections[connId]
	if !ok {
		return
	}

	// Cancel subscribe
	a.onCancelSub <- a.connectionsSub[connId]
	// Cancel subscribe in the handler
	a.handler(conn).OnLeave(conn)
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

func (a *agent) handler(conn frontier.Conn) *handler {
	index := conn.GetId() % a.parallel
	return a.handlers[index]
}
