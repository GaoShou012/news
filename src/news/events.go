package news

import proto_news "im/proto/news"

type EventOnSubscribe struct {
	Client   *Client
	Channels []string
}

type EventOnLeave struct {
	Client *Client
}

type EventOnNews struct {
	Item *proto_news.NewsItem
}
