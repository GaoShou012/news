package news

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"im/src/frontier"
	"im/src/im"
)

type handler struct {
	newsCache        chan *Item
	clients          map[int]*Client
	clientsOnSubList map[string]SubClients

	connOnSwitchToDirectMode chan frontier.Conn
	addConn                  chan frontier.Conn
	delConn                  chan frontier.Conn

	uploadSubList chan *EventUploadSubscribe

	task chan *Task
}

func (h *handler) init() {
	h.newsCache = make(chan *Item, 100000)
	h.clients = make(map[int]*Client)
	h.clientsOnSubList = make(map[string]SubClients)

	h.connOnSwitchToDirectMode = make(chan frontier.Conn, 100000)
	h.addConn = make(chan frontier.Conn, 100000)
	h.delConn = make(chan frontier.Conn, 100000)
	h.uploadSubList = make(chan *EventUploadSubscribe, 100000)

	go func() {
		for {
			select {
			// 推送给指定终端
			case task := <-h.task:
			// 广播给指定订阅组
			case conn := <-h.addConn:
				// 连接接入
				id := conn.GetId()
				_, ok := h.clients[id]
				if ok {
					glog.Errorln(fmt.Sprintf("连接接入：连接ID已存在"))
					break
				}
				cli := &Client{
					Conn:          conn,
					IsCaching:     true,
					NewsCache:     list.New(),
					LastMessageId: make(map[string]string),
				}
				h.clients[id] = cli
			case conn := <-h.delConn:
				// 连接断开
				id := conn.GetId()
				cli, ok := h.clients[id]
				if !ok {
					glog.Errorln(fmt.Sprintf("连接断开：连接ID不存在"))
					break
				}
				delete(h.clients, id)
				h.delClient(cli)
			case event := <-h.uploadSubList:
				// 客户端上传监听列表
				conn, message := event.Conn, event.Message
				cli, ok := h.clients[event.Conn.GetId()]
				if !ok {
					message.Head.Status = 1
					message.Body = fmt.Sprintf("Cli不存在")
					msg, err := im.Encode(message)
					if err != nil {
						glog.Errorln(err)
						break
					}
					conn.Sender(msg)
					break
				}
				subscribeList, ok := message.Body.([]string)
				if !ok {
					message.Head.Status = 1
					message.Body = fmt.Sprintf("监听列表转换失败")
					msg, err := im.Encode(message)
					if err != nil {
						glog.Errorln(err)
						break
					}
					conn.Sender(msg)
					break
				}
				h.delClient(cli)
				h.addSubList(cli, subscribeList)
			case item := <-h.newsCache:
				// 新消息进行推送
				sub, ok := h.clientsOnSubList[item.Key]
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
			case conn := <-h.connOnSwitchToDirectMode:
				cli, ok := h.clients[conn.GetId()]
				if !ok {
					break
				}
				for {
					item := cli.NewsCachePop()
					if item == nil {
						break
					}
					if cli.IsNewMessageId(item.Key, item.Id) {
						j, err := json.Marshal(item)
						if err != nil {
							glog.Errorln(err)
							continue
						}
						cli.Conn.Sender(j)
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
		sub, ok := h.clientsOnSubList[key]
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
		sub, ok := h.clientsOnSubList[key]
		if !ok {
			continue
		}
		delete(sub, id)
	}
}

// 移除客户端，清空客户端已经监听的所有列表
func (h *handler) delClient(cli *Client) {
	id := cli.Conn.GetId()
	for _, sub := range h.clientsOnSubList {
		delete(sub, id)
	}
}
func (h *handler) AddConn(conn frontier.Conn) {
	h.addConn <- conn
}
func (h *handler) DelConn(conn frontier.Conn) {
	h.delConn <- conn
}
func (h *handler) UploadSubList(conn frontier.Conn, message *im.Message) {
	h.uploadSubList <- &EventUploadSubscribe{
		Conn:    conn,
		Message: message,
	}
}
func (h *handler) ConnSwitchToDirectMode(conn frontier.Conn) {
	h.connOnSwitchToDirectMode <- conn
}
