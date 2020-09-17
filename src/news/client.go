package news

import (
	"github.com/GaoShou012/frontier"
	proto_message "im/proto/message"
	proto_news "im/proto/news"
)

type Client struct {
	Conn      frontier.Conn
	IsCaching bool
	News      map[string]*proto_news.NewsItem
}

func (c *Client) IsNewMessageId(message *proto_news.NewsItem) bool {
	item, ok := c.News[message.Key]
	if !ok {
		return true
	}
	return proto_message.IsNewMessageId(item.Id, message.Id)
}

func (c *Client) Sender(message []byte) {
	c.Conn.Sender(message)
}
