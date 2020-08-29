package news

import (
	proto_news "im/proto/news"
	"im/src/im"
)

type Message struct {
	Type string
	Key  string
	Val  string
}

type Item struct {
	Id  string
	Key string
	Val string
}

func ResponseNewsItems(items []*proto_news.NewsItem) *im.Message {
	msg := &im.Message{
		Head: &im.Head{
			Id:     0,
			Status: 0,
			Api:    "",
		},
		Body: items,
	}
	return msg
}

func ResponseError(status int, err error) *im.Message {
	msg := &im.Message{
		Head: &im.Head{
			Id:     0,
			Status: status,
			Api:    "error",
		},
		Body: err,
	}
	return msg
}
