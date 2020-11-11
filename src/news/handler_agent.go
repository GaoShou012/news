package news

import (
	"fmt"
	"github.com/GaoShou012/frontier"
	"github.com/golang/glog"
	proto_news "im/proto/news"
	"im/utils/logger"
	"runtime"
	"sync"
	"time"
)

type HandlerAgent struct {
	clients      map[int]*Client
	clientsMutex sync.Mutex
	parallel     int
	handlers     []*Handler

	puller struct {
		channels []string
		mutex    sync.Mutex
		bufFlag  int
		buf      []map[string]bool
	}
}

func (h *HandlerAgent) Init() {
	h.clients = make(map[int]*Client)
	h.parallel = runtime.NumCPU()
	h.handlers = make([]*Handler, h.parallel)
	for i := 0; i < h.parallel; i++ {
		handler := &Handler{}
		handler.OnInit()
		h.handlers[i] = handler
	}

	// init puller
	h.puller.buf = make([]map[string]bool, 2)
	for i := 0; i < 2; i++ {
		h.puller.buf[i] = make(map[string]bool)
	}
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			<-ticker.C
			h.puller.mutex.Lock()
			bufFlag := h.puller.bufFlag
			if h.puller.bufFlag == 0 {
				h.puller.bufFlag = 1
			} else {
				h.puller.bufFlag = 0
			}
			h.puller.mutex.Unlock()

			channels := h.puller.buf[bufFlag]
			h.puller.buf[bufFlag] = make(map[string]bool)
			for channelName, _ := range channels {
				// TODO 拉取频道的消息
				fmt.Println("拉取频道：", channelName)
				news := &proto_news.NewsItem{
					Key: "abc",
					Val: []byte("aaa"),
					Id:  "21313",
				}
				h.OnNews(news)
			}
		}
	}()
}

func (h *HandlerAgent) index(id int) int {
	return id % h.parallel
}

func (h *HandlerAgent) addChannelToPuller(channel string) {
	h.puller.mutex.Lock()
	defer h.puller.mutex.Unlock()
	h.puller.buf[h.puller.bufFlag][channel] = true
}

/*
	客户端订阅频道
*/
func (h *HandlerAgent) subscribe(conn frontier.Conn, sub *proto_news.Subscribe) {
	h.clientsMutex.Lock()
	defer h.clientsMutex.Unlock()

	logger.Compare(logger.LvInfo, "subscribe", conn)
	client := NewClient(conn)
	h.clients[conn.GetId()] = client
	index := h.index(conn.GetId())
	h.handlers[index].Subscribe(client, sub.Channels)
}

/*
	客户端取消订阅
*/
func (h *HandlerAgent) unsubscribe(conn frontier.Conn) {
	h.clientsMutex.Lock()
	defer h.clientsMutex.Unlock()

	logger.Compare(logger.LvInfo, "unsubscribe", conn)
	client, ok := h.clients[conn.GetId()]
	if !ok {
		return
	}
	delete(h.clients, conn.GetId())
	index := h.index(conn.GetId())
	h.handlers[index].Unsubscribe(client)
}

/*
	客户端消息
	接受客户端的，订阅 & 取消订阅 操作
*/
func (h *HandlerAgent) OnMessage(conn frontier.Conn, message []byte) {
	// 解码消息
	msg := &proto_news.Message{}
	{
		err := proto_news.Decode(message, msg)
		if err != nil {
			glog.Errorln(err)
			return
		}
	}

	// 消息处理
	switch msg.Head.Path {
	case "subscribe":
		sub := &proto_news.Subscribe{}
		if err := proto_news.Decode(msg.Body, sub); err != nil {
			glog.Errorln(err)
			return
		}
		h.subscribe(conn, sub)
		break
	default:
		glog.Errorln("Unknown message type", msg.Head.Path)
		break
	}
}

func (h *HandlerAgent) OnNews(news *proto_news.NewsItem) {
	for i := 0; i < h.parallel; i++ {
		h.handlers[i].onNews <- news
	}
}

func (h *HandlerAgent) OnLeave(conn frontier.Conn) {
	h.unsubscribe(conn)
}
