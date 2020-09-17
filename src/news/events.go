package news

import (
	"github.com/GaoShou012/frontier"
	proto_news "im/proto/news"
)

type EventOnSubscribe struct {
	Conn    frontier.Conn
	Message *proto_news.Subscribe
}
