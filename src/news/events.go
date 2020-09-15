package news

import (
	proto_news "im/proto/news"
	"im/src/frontier"
)

type EventOnSubscribe struct {
	Conn    frontier.Conn
	Message *proto_news.Subscribe
}
