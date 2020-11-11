package news

import proto_news "im/proto/news"

type JobNewsItemStreamEncode struct {
	input  *proto_news.NewsItem
	output []byte
	err    chan error
}
type JobNewsItemStreamDecode struct {
	input  []byte
	output *proto_news.NewsItem
	err    chan error
}
