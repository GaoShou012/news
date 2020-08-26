package news

import (
	"container/list"
	"encoding/json"
	"github.com/prometheus/common/log"
)

var News news
var Handler handler

func init() {
	News.OnInit()
}

type handler struct {
	clients list.List

	addClient                   chan *Client
	delClient                   chan *Client
	clientsOnSwitchToDirectMode chan *Client
	newsCache                   chan *Item
}

func (h *handler) OnInit() {
	h.clientsOnSwitchToDirectMode = make(chan *Client, 100000)
	h.newsCache = make(chan *Item, 100000)
	go func() {
		for {
			select {
			case item := <-h.newsCache:
			case cli := h.addClient:
					
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
