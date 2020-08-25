package news

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"runtime"
	"wchatv1/src/frontier"
	"wchatv1/src/im"
)

type news struct {
	// 所有参与订阅的client，订阅列表
	clients map[int]map[string]bool
	// 所有频道的clients列表
	clientsSubList   []map[string]frontier.Conn
	clientsOnSubList []map[string]*list.List
	caches           []chan *Message
	events           []chan *Event
}

func (n *news) OnInit() {
	n.clients = make(map[int]map[string]bool)
	n.clientsSubList = make([]map[string]frontier.Conn, runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		n.clientsSubList[i] = make(map[string]frontier.Conn)
	}

	n.caches = make([]chan *Message, runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		n.caches[i] = make(chan *Message)
	}

	// 处理新消息的推送
	for i := 0; i < runtime.NumCPU(); i++ {
		cache := n.caches[i]
		clientsOnSubList := n.clientsOnSubList[i]
		events := n.events[i]
		go func(cache chan *Message, clientsOnSubList map[string]*list.List) {
			for {
				select {
				case message := <-cache:
					clients, ok := clientsOnSubList[message.Kind]
					if !ok {
						continue
					}
					j, err := json.Marshal(message)
					if err != nil {
						glog.Errorln(err)
						continue
					}
					for ele := clients.Front(); ele != nil; ele = ele.Next() {
						conn := ele.Value.(frontier.Conn)
						conn.Writer(j)
					}
				case event := <-
				}
			}
		}(cache, clientsOnSubList,events)
	}
}

func (n *news) OnMessage(conn frontier.Conn, message *im.Message) {
	switch message.Head.BusinessApi {
	case "subscribe.list.upload":
		m, ok := message.Body.(map[string]map[string]bool)
		if !ok {
			im.ResponseError(conn, message, fmt.Errorf("订阅列表格式错误"))
			return
		}
		// 检验订阅列表是否有效
		break
	case "subscribe.list.download":
		break
	default:
		im.ResponseError(conn, message, fmt.Errorf("business api is not existsing"))
	}
}
