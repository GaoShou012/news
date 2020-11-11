package news

import (
	proto_news "im/proto/news"
	"im/utils/logger"
)

type Handler struct {
	onSubscribe chan *EventOnSubscribe
	onLeave     chan *EventOnLeave
	onNews      chan *proto_news.NewsItem
	channels    map[ChannelName]*Channel
}

func (h *Handler) OnInit() {
	h.onSubscribe = make(chan *EventOnSubscribe, 100000)
	h.onLeave = make(chan *EventOnLeave, 100000)
	h.onNews = make(chan *proto_news.NewsItem, 100000)
	h.channels = make(map[ChannelName]*Channel)
	h.handler()
}

/*
	添加订阅
	订阅后的频道锚，保存在client
*/
func (h *Handler) subscribe(cli *Client, channels []string) {
	for _, channelName := range channels {
		chName := ChannelName(channelName)
		channel, ok := h.channels[chName]
		if !ok {
			channel = &Channel{}
			h.channels[chName] = channel
		}
		cli.ChannelsAnchors[chName] = channel.Add(cli)
	}
}

/*
	取消订阅
	遍历client的所有频道锚点，进行取消订阅
*/
func (h *Handler) unsubscribe(cli *Client) []ChannelName {
	var channels []ChannelName

	for chName, anchor := range cli.ChannelsAnchors {
		channel, ok := h.channels[chName]
		if !ok {
			logger.Compare(logger.LvError, "channel is not exists")
			continue
		}

		// to delete the anchor of the client
		channel.Del(anchor)

		// to delete the channel
		if channel.Len() == 0 {
			delete(h.channels, chName)
		}

		// return the unsubscribe channels
		channels = append(channels, chName)
	}
	cli.CleanNewsRecord()
	cli.CleanChannelsAnchors()

	return channels
}

func (h *Handler) handler() {
	go func() {
		for {
			select {
			case event := <-h.onSubscribe:
				// 取消之前的订阅
				h.unsubscribe(event.Client)
				// 重新订阅
				h.subscribe(event.Client, event.Channels)

				logger.Compare(logger.LvInfo, "添加订阅", event.Channels)
			case event := <-h.onLeave:
				// 取消之间的订阅
				channels := h.unsubscribe(event.Client)

				logger.Compare(logger.LvInfo, "取消订阅", channels)
			case news := <-h.onNews:
				logger.Compare(logger.LvInfo, "收到新消息", news)

				// 检查频道是否有客户订阅
				chName := ChannelName(news.Key)
				channel, ok := h.channels[chName]
				if !ok {
					continue
				}

				// 序列化数据
				data, err := DefaultCodec.NewsItemStreamEncode(news)
				if err != nil {
					logger.Compare(logger.LvError, err)
					continue
				}

				// Info log
				logger.Compare(logger.LvInfo, "推送消息", news)

				// 判断消息是否最新（相对客户端而言）
				// 最新的消息，推送给客户端，并且记录当前消息
				for ele := channel.Clients.Front(); ele != nil; ele = ele.Next() {
					cli, ok := ele.Value.(*Client)
					if !ok {
						logger.Compare(logger.LvError, "断言客户类型失败")
						continue
					}

					logger.Compare(logger.LvInfo, "publish cli", cli)

					if cli.IsNewMessageId(news) == false {
						continue
					}

					logger.Compare(logger.LvInfo, "publish cli news", news)

					cli.Conn.Sender(data)
					cli.News[chName] = news
				}
			}
		}
	}()
}

/*
	订阅
*/
func (h *Handler) Subscribe(cli *Client, channels []string) {
	event := &EventOnSubscribe{
		Client:   cli,
		Channels: channels,
	}
	h.onSubscribe <- event
}

/*
	取消订阅
*/
func (h *Handler) Unsubscribe(cli *Client) {
	event := &EventOnLeave{Client: cli}
	h.onLeave <- event
}

/*
	推送新消息
*/
func (h *Handler) OnNews(news *proto_news.NewsItem) {
	h.onNews <- news
}
