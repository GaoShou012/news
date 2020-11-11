package news

import (
	"fmt"
	"github.com/golang/glog"
	proto_news "im/proto/news"
	"testing"
)

func TestChannelNews_Set(t *testing.T) {
	fmt.Println("测试频道news set")

	item := &proto_news.NewsItem{
		Key: "testing",
		Val: "abc",
		Id:  "",
	}

	channelNews := &ChannelNews{RedisClient: newRedis()}
	str, err := channelNews.Set(item)
	if err != nil {
		glog.Errorln(err)
		return
	}
	fmt.Println("item id:", str)
}

func TestChannelNews_Get(t *testing.T) {
	fmt.Println("测试频道news get")

	channelNews := &ChannelNews{RedisClient: newRedis()}
	item, err := channelNews.Get("testing")
	if err != nil {
		glog.Errorln(err)
		return
	}
	fmt.Println("item", item)
}
