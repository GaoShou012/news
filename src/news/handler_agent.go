package news

import (
	"context"
	"github.com/golang/glog"
	"im/config"
	proto_news "im/proto/news"
	"im/src/frontier"
	"im/src/im"
	"im/src/logger"
	"im/src/parallel"
	"runtime"
	"sync"
)

var LogLevel int

type handlerAgent struct {
	mutex   sync.Mutex
	clients map[int]*Client
	// 上传订阅列表
	uploadSubList []chan *UploadSubList

	// 所有频道的消息缓存
	channelsNewsCache sync.Map

	// 频道的订阅者
	channels []map[string]ChannelSubscribe
	// 推送消息给指定客户端
	pubNewsToClient []chan *PubNewsToClient
	// 推送消息给订阅组
	pubNews []chan *proto_news.NewsItem

	// 拉取消息
	pullNews chan *PullNewsTask

	parallelForUploadSubList *parallel.Parallel
}

func (h *handlerAgent) init() {
	for i := 0; i < runtime.NumCPU(); i++ {
		go func(i int) {
			channels := h.channels[i]
			for {
				select {
				case pub := <-h.pubNewsToClient[i]:
					// 推送消息给指定客户端

					client, items := pub.Client, pub.Items
					var newItems []*proto_news.NewsItem
					for _, item := range items {
						if client.Channels.IsSub(item.Key) == false {
							continue
						}
						if client.Channels.IsNew(item) == false {
							continue
						}
						client.Channels.SetItem(item)
						newItems = append(newItems, item)
					}
					if len(newItems) <= 0 {
						break
					}

					// 生成消息报文
					message := ResponseNewsItems(newItems)
					msg, err := im.Encode(message)
					if err != nil && LogLevel >= logger.LogLevelError {
						glog.Errorln(err)
						break
					}
					// 推送消息
					client.Conn.Sender(msg)
				case item := <-h.pubNews[i]:
					// 推送消息给订阅者

					// 生成消息报文
					message := ResponseNewsItems([]*proto_news.NewsItem{item})
					msg, err := im.Encode(message)
					if err != nil && LogLevel >= logger.LogLevelError {
						glog.Errorln(err)
						break
					}
					// 推送消息至订阅此频道的所有客户端
					channelSubscribe, ok := channels[item.Key]
					if !ok {
						break
					}
					for _, client := range channelSubscribe {
						if client.Channels.IsNew(item) == false {
							continue
						}
						client.Channels.SetItem(item)
						client.Conn.Sender(msg)
					}
				case event := <-h.uploadSubList[i]:
					// 上传订阅列表

					client, subList := event.Client, event.SubList

					client.Channels = nil
					client.Channels = &Channels{}
					for _, channelName := range subList {
						client.Channels.Sub(channelName)
					}

					// 取消订阅
					for _, channelSubscribe := range channels {
						channelSubscribe.DelClient(client)
					}

					// 订阅
					for _, channelName := range subList {
						channels[channelName].AddClient(client)
					}

				}
			}
		}(i)
	}
}

// 检查本地频道缓存是否存在，如果不存在从服务拉取
// 推送消息给指定的客户端
func (h *handlerAgent) puller() {
	h.pullNews = make(chan *PullNewsTask, 100000)
	for {
		task := <-h.pullNews
		client, subList := task.Client, task.SubList
		var items []*proto_news.NewsItem
		var channels []string
		for _, key := range subList {
			val, ok := h.channelsNewsCache.Load(key)
			if !ok {
				channels = append(channels, key)
				continue
			}
			item := val.(*proto_news.NewsItem)
			items = append(items, item)
		}
		if len(items) > 0 {
			h.pubNewsToClient[h.parallel(client.Conn.GetId())] <- &PubNewsToClient{
				Client: client,
				Items:  items,
			}
		}
		if len(channels) == 0 {
			continue
		}

		req := &proto_news.GetNewsReq{
			Keys:  channels,
			Count: 1024,
		}
		rsp, err := config.NewsServiceConfig.ServiceClient().GetNews(context.TODO(), req)
		if err != nil && LogLevel >= logger.LogLevelError {
			glog.Errorln(err)
			continue
		}
		if len(rsp.Items) == 0 {
			continue
		}
		h.pubNewsToClient[h.parallel(client.Conn.GetId())] <- &PubNewsToClient{
			Client: client,
			Items:  rsp.Items,
		}
	}
}

func (h handlerAgent) parallel(id int) int {
	return id % runtime.NumCPU()
}
func (h *handlerAgent) AddConn(conn frontier.Conn) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	_, ok := h.clients[conn.GetId()]
	if ok {
		if LogLevel >= logger.LogLevelError {
			glog.Errorln("Conn ID已经存在\n")
			return
		}
	}
	client := &Client{
		Conn:      nil,
		IsCaching: false,
		NewsCache: nil,
		Channels:  nil,
	}
	h.clients[conn.GetId()] = client
}
func (h *handlerAgent) DelConn(conn frontier.Conn) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	_, ok := h.clients[conn.GetId()]
	if !ok {
		if LogLevel >= logger.LogLevelError {
			glog.Errorln("Conn ID不存在\n")
			return
		}
	}
	delete(h.clients, conn.GetId())
}
func (h *handlerAgent) UploadSubscribeList(conn frontier.Conn, message *im.Message) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	cli, ok := h.clients[conn.GetId()]
	if !ok {
		if LogLevel >= logger.LogLevelError {
			glog.Errorln("Conn ID不存在\n")
			return
		}
	}

	go func(cli *Client, message *im.Message) {
		h.parallelForUploadSubList.Do()
		defer h.parallelForUploadSubList.Done()

		subList, ok := message.Body.([]string)
		if !ok {
			message.Head.Status = 1
			message.Body = "订阅列表转换失败"
			msg, err := im.Encode(message)
			if err != nil {
				glog.Errorln(err)
				return
			}
			cli.Conn.Sender(msg)
		}
		if len(subList) <= 0 || len(subList) >= 255 {
			message.Head.Status = 1
			message.Body = "订阅列表数量超过限制"
			msg, err := im.Encode(message)
			if err != nil {
				glog.Errorln(err)
				return
			}
			cli.Conn.Sender(msg)
		}

		index := h.parallel(cli.Conn.GetId())
		h.uploadSubList[index] <- &UploadSubList{
			Client:  cli,
			SubList: subList,
		}
		h.pullNews <- &PullNewsTask{
			Client:  cli,
			SubList: subList,
		}
	}(cli, message)
}
