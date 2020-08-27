package news

import (
	"container/list"
	"github.com/prometheus/common/log"
	"im/src/frontier"
)

type handler struct {
	clients map[int]*Client

	clientsOnSwitchToDirectMode chan *Client
	newsCache                   chan *Item

	addConn chan frontier.Conn
	delConn chan frontier.Conn

	uploadSubList chan *EventUploadSubscribe
	deleteClient  chan *Client

	eventAddSubList chan *EventAddSubList
	eventDelSubList chan *EventDelSubList

	subList map[string]SubClients
}

func (h *handler) init() {
	h.clientsOnSwitchToDirectMode = make(chan *Client, 100000)
	h.newsCache = make(chan *Item, 100000)
	go func() {
		for {
			select {
			case item := <-h.newsCache:
				sub, ok := h.subList[item.Key]
				if !ok {
					break
				}
				msg, err := json.Marshal(item)
				if err != nil {
					break
				}
				for _, cli := range sub {
					if cli.IsCaching {
						cli.NewsCachePush(item)
					} else {
						cli.Conn.Sender(msg)
					}
				}
			case conn := <-h.addConn:
				id := conn.GetId()
				_, ok := h.clients[id]
				if ok {
					return
				}
				cli := &Client{
					Conn:          conn,
					IsCaching:     true,
					NewsCache:     list.New(),
					LastMessageId: make(map[string]string),
				}

			case conn := <-h.delConn:
			case event := <-h.eventAddSubList:
				h.addSubList(event.Client, event.SubscribeList)
			case event := <-h.eventDelSubList:
				h.delSubList(event.Client, event.SubscribeList)
			case event := <-h.uploadSubList:
				h.delClient(event.Client)
				h.addSubList(event.Client, event.SubscribeList)
			case cli := <-h.deleteClient:
				h.delClient(cli)

			case cli := <-h.clientsOnSwitchToDirectMode:
				for {
					item := cli.NewsCachePop()
					if item == nil {
						break
					}
					if cli.IsNewMessageId(item.Key, item.Id) {
						j, err := json.Marshal(item)
						if err != nil {
							log.Errorln(err)
							continue
						}
						cli.Conn.Writer(j)
					}
				}
				cli.NewsCache = nil
				cli.IsCaching = false
			}
		}
	}()
}

// 添加客户端需要监听的列表
func (h *handler) addSubList(cli *Client, subList []string) {
	id := cli.Conn.GetId()
	for _, key := range subList {
		sub, ok := h.subList[key]
		if !ok {
			sub = make(SubClients)
		}
		sub[id] = cli
	}
}

// 删除客户端需要清除的监听列表
func (h *handler) delSubList(cli *Client, subList []string) {
	id := cli.Conn.GetId()
	for _, key := range subList {
		sub, ok := h.subList[key]
		if !ok {
			continue
		}
		delete(sub, id)
	}
}

// 移除客户端，清空客户端已经监听的所有列表
func (h *handler) delClient(cli *Client) {
	id := cli.Conn.GetId()
	for _, sub := range h.subList {
		delete(sub, id)
	}
}

func (h *handler) AddConn(conn frontier.Conn) {

}
func (h *handler) DelConn(conn frontier.Conn) {

}

func (h *handler) UploadSubList(cli *Client, subList []string) {
	h.uploadSubList <- &EventUploadSubscribe{
		Client:        cli,
		SubscribeList: subList,
	}
}
func (h *handler) CleanSubListByClient(cli *Client) {
	h.deleteClient <- cli
}
func (h *handler) ClientToDirectMode(cli *Client) {
	h.clientsOnSwitchToDirectMode <- cli
}
