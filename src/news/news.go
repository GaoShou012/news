package news

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"im/src/frontier"
	"im/src/im"
	"runtime"
)

type news struct {
	// 所有频道的conn列表
	subList []map[string]*ConnList
	// 消息缓存
	caches []chan *Message

	eventOnCleanSubscribe  []chan *EventCleanSubscribe
	eventOnUploadSubscribe []chan *EventUploadSubscribe

	clientsListOnSyncNews *list.List
}

func (n *news) OnInit() {
	for i := 0; i < runtime.NumCPU(); i++ {
		go func(i int) {
			subList := n.subList[i]
			for {
				select {
				case message := <-n.caches[i]:
					// 处理新消息的推送
					connList, ok := subList[message.Type]
					if !ok {
						continue
					}
					j, err := json.Marshal(message)
					if err != nil {
						glog.Errorln(err)
						continue
					}
					for ele := connList.GetConnections().Front(); ele != nil; ele = ele.Next() {
						conn := ele.Value.(frontier.Conn)
						conn.Writer(j)
					}
				case event := <-n.eventOnUploadSubscribe[i]:
					// 上传订阅列表
					conn := event.Conn
					for _, v := range event.SubscribeList {
						if _, ok := subList[v]; !ok {
							connList := &ConnList{}
							connList.Init()
							subList[v] = connList
						}
						subList[v].Push(conn)
					}
				case event := <-n.eventOnCleanSubscribe[i]:
					// 移除连接的订阅
					conn := event.Conn
					for _, connList := range subList {
						connList.Remove(conn)
					}
				}
			}
		}(i)
	}
}

func (n *news) OnOpen(conn frontier.Conn) {
	HandlerAgent.AddConn(conn)
}
func (n *news) OnClose(conn frontier.Conn) {
	HandlerAgent.DelConn(conn)
}
func (n *news) OnMessage(conn frontier.Conn, message *im.Message) {
	switch message.Head.BusinessApi {
	case "subscribe.list.upload":
		m, ok := message.Body.([]string)
		if !ok {
			im.ResponseError(conn, message, fmt.Errorf("订阅列表格式错误"))
			return
		}
		if len(m) > 255 {
			im.ResponseError(conn,message,fmt.Errorf(""))
		}
		// 检验订阅列表是否有效

		// 添加到处理程序
		HandlerAgent.UploadSubscribeList(conn, message)
		break
	case "subscribe.list.download":
		break
	default:
		im.ResponseError(conn, message, fmt.Errorf("business api is not existsing"))
	}
}

func (n *news) handler() {
	n.clientsListOnSyncNews = list.New()

	// 同步最新的消息
	for i := 0; i < 100; i++ {
		go func() {
			for {
				ele := n.clientsListOnSyncNews.Front()
				if ele == nil {
					continue
				}
			}
		}()
	}

	for i := 0; i < 100; i++ {
		go func() {
			for {
				cli := n.ta
			}
		}()
	}

}
