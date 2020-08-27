package news

import (
	"container/list"
	"encoding/json"
	"errors"
	"github.com/prometheus/common/log"
	"im/src/frontier"
)

var News news
var Handler handler

func init() {
	News.OnInit()
}

type handler struct {
	clientsAnchor map[int]*list.Element
	clients       list.List
	
	addClient                   chan *Client
	delClient                   chan *Client
	clientsOnSwitchToDirectMode chan *Client
	newsCache                   chan *Item

	uploadSubList   chan *EventUploadSubscribe
	downloadSubList chan *EventDownloadSubscribe
	cleanSubList    chan *EventCleanSubscribe

	subList map[string]map[int]*Client
}

func (h *handler) getClient(id int) *Client {
	ele, ok := h.clientsAnchor[id]
	if !ok {
		return nil
	}
	cli := ele.Value.(*Client)
	return cli
}
func (h *handler) setClient(conn frontier.Conn) error {
	id := conn.GetId()
	_,ok := h.clientsAnchor[id]
	if ok {
		return errors.New("客户端已经存在")
	}

}

func (h *handler) OnInit() {
	h.clientsOnSwitchToDirectMode = make(chan *Client, 100000)
	h.newsCache = make(chan *Item, 100000)
	go func() {
		for {
			select {
			case item := <-h.newsCache:
				// 新消息推送
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
			case event := <-h.uploadSubList:
				for _,v := range event.SubscribeList {
					sub,ok := h.subList[v]
					if !ok {
						h.subList[v] = make(map[int]*)
					}
				}

			case cli := <-h.addClient:
				// 新客户端
				for _, v := range cli.SubList {
					sub, ok := h.subList[v]
					if !ok {
						h.subList[v] = make(map[int]*Client)
						sub = h.subList[v]
					}
					id := cli.Conn.GetId()
					sub[id] = cli
				}
			case cli := <-h.delClient:
				id := cli.Conn.GetId()
				for _, sub := range h.subList {
					delete(sub, id)
				}
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

func (h *handler) ClientToDirectMode(cli *Client) {
	h.clientsOnSwitchToDirectMode <- cli
}
