package news

import (
	"container/list"
	"context"
	"encoding/json"
	"github.com/golang/glog"
	"im/config"
	proto_news "im/proto/news"
	"time"
)


/*
	把订阅列表，同步到微服务
	并且通过kafka返回响应的消息
*/


type Sync struct {
	clients *list.List
	task    chan *Client
}

func (s *Sync) Push(cli *Client) {
	s.clients.PushBack(cli)
}
func (s *Sync) Pop() *Client {
	ele := s.clients.Front()
	if ele == nil {
		return nil
	}
	s.clients.Remove(ele)
	return ele.Value.(*Client)
}

func (s *Sync) Handler() {
	go func() {
		ticker := time.NewTicker(time.Millisecond * 100)
		for {
			<-ticker.C
			cli := s.Pop()
			if cli == nil {
				continue
			}
			if cli.Conn.IsWorking() == false {
				continue
			}
			s.task <- cli
		}
	}()
	for i := 0; i < 100; i++ {
		go func() {
			for {
				cli := <-s.task
				req := &proto_news.GetNewsReq{Keys: cli.SubList, Count: 50}
				rsp, err := config.NewsServiceConfig.ServiceClient().GetNews(context.TODO(), req)
				if err != nil {
					glog.Errorln(err)
					continue
				}

				var items []Item
				for _, v := range rsp.Items {
					cli.LastMessageId[v.Key] = v.Id
					items = append(items, Item{
						Key: v.Key,
						Val: v.Val,
					})
				}
				if len(items) == 0 {
					continue
				}

				msg := ResponseNewsItems(items)
				j, err := json.Marshal(msg)
				cli.Conn.Writer(j)

				if len(items) < 50 {
					Handler.ClientToDirectMode(cli)
					continue
				}

				s.Push(cli)
			}
		}()
	}

}
