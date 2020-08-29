package news

import (
	proto_news "im/proto/news"
	"im/src/im"
)

type ChannelSubscribe map[int]*Client

func (c ChannelSubscribe) AddClient(cli *Client) {
	c[cli.Conn.GetId()] = cli
}
func (c ChannelSubscribe) DelClient(cli *Client) {
	delete(c, cli.Conn.GetId())
}

type Channels map[string]*proto_news.NewsItem

func (c Channels) IsSub(key string) bool {
	_, ok := c[key]
	return ok
}
func (c Channels) Sub(key string) {
	c[key] = nil
}
func (c Channels) IsNew(newItem *proto_news.NewsItem) bool {
	oldItem, ok := c[newItem.Key]
	if !ok {
		return true
	}
	if oldItem == nil {
		return true
	}
	return im.IsNewMessageId(oldItem.Id, newItem.Id)
}
func (c Channels) SetItem(item *proto_news.NewsItem) {
	c[item.Key] = item
}
func (c Channels) GetItem(key string) (*proto_news.NewsItem, bool) {
	item, ok := c[key]
	return item, ok
}
