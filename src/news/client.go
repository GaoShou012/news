package news

import (
	"container/list"
	"github.com/GaoShou012/frontier"
	proto_im "im/proto/im"
	proto_news "im/proto/news"
)

func NewClient(conn frontier.Conn) *Client {
	return &Client{
		Conn:            conn,
		News:            make(map[ChannelName]*proto_news.NewsItem),
		ChannelsAnchors: make(map[ChannelName]*list.Element),
	}
}

/*
	Conn ws连接
	News 记录每个频道下已经收到的最新内容
	ChannelsAnchors 记录每个频道下，自己处于频道内的锚点，用于快于索引自己的位置，便于从频道中移除
*/
type Client struct {
	Conn            frontier.Conn
	News            map[ChannelName]*proto_news.NewsItem
	ChannelsAnchors map[ChannelName]*list.Element
}

/*
	获取频道，已经推送给客户端的最后消息ID
	如果没有进行推送，返回空字符串
*/
func (c *Client) GetLastMessageId(channelName ChannelName) string {
	news, ok := c.News[channelName]
	if !ok {
		return ""
	} else {
		return news.Id
	}
}

/*
	判断频道消息，是否最新
*/
func (c *Client) IsNewMessageId(newsItem *proto_news.NewsItem) bool {
	oldMsgId := c.GetLastMessageId(ChannelName(newsItem.Key))
	if oldMsgId == "" {
		return true
	}
	return proto_im.IsNewMessageId(oldMsgId, newsItem.Id)
}

/*
	清空客户端订阅的频道的描点
*/
func (c *Client) CleanChannelsAnchors() {
	c.ChannelsAnchors = nil
	c.ChannelsAnchors = make(map[ChannelName]*list.Element)
}

/*
	清空客户端各个频道下的消息
*/
func (c *Client) CleanNewsRecord() {
	c.News = nil
	c.News = make(map[ChannelName]*proto_news.NewsItem)
}
