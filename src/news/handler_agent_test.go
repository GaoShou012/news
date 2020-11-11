package news

import (
	"fmt"
	proto_news "im/proto/news"
	"testing"
	"time"
)

func TestHandlerAgent_Init(t *testing.T) {
	fmt.Println("测试agent init")

	agent := &HandlerAgent{}
	agent.Init()
	agent.addChannelToPuller("abc")
}

func TestHandlerAgent_OnNews(t *testing.T) {
	fmt.Println("测试agent OnNews")

	news := &proto_news.NewsItem{
		Key: "abc",
		Val: "abc",
		Id:  "213",
	}

	agent := &HandlerAgent{}
	agent.Init()
	agent.OnNews(news)

	time.Sleep(time.Second * 10)
}
