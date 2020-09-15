package news

import (
	proto_news "im/proto/news"
	"im/src/frontier"
	"im/src/im"
)

type UploadSubList struct {
	Client  *Client
	SubList []string
}

type EventOnSubscribe struct {
	Conn    frontier.Conn
	Message *proto_news.Subscribe
}

type EventUploadSubscribe struct {
	Conn          frontier.Conn
	Message       *im.Message
	SubscribeList []string
}

type EventCleanSubscribe struct {
	Client *Client
}

type EventDownloadSubscribe struct {
	Conn frontier.Conn
}

type EventAddSubList struct {
	Client        *Client
	SubscribeList []string
}
type EventDelSubList struct {
	Client        *Client
	SubscribeList []string
}
