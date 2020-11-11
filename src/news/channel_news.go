package news

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/golang/glog"
	"google.golang.org/protobuf/proto"
	proto_news "im/proto/news"
)

func KeyOfChannelNewsStream(key string) string {
	return fmt.Sprintf("im:news:%s", key)
}

func EncodeForStream(item *proto_news.NewsItem) ([]byte, error) {
	data, err := proto.Marshal(item)
	return data, err
}

func DecodeForStream(data []byte) (*proto_news.NewsItem, error) {
	item := &proto_news.NewsItem{}
	err := proto.Unmarshal(data, item)
	return item, err
}

type ChannelNews struct {
	RedisClient *redis.ClusterClient
}

func (n *ChannelNews) Set(item *proto_news.NewsItem) (string, error) {
	data, err := EncodeForStream(item)
	if err != nil {
		return "", err
	}
	stream := KeyOfChannelNewsStream(item.Key)
	values := make(map[string]interface{})
	values["payload"] = data
	xAddArgs := &redis.XAddArgs{
		Stream:       stream,
		MaxLen:       5,
		MaxLenApprox: 0,
		ID:           "*",
		Values:       values,
	}
	str, err := n.RedisClient.XAdd(xAddArgs).Result()
	return str, err
}

func (n *ChannelNews) Get(key string) (*proto_news.NewsItem, error) {
	stream := KeyOfChannelNewsStream(key)
	res, err := n.RedisClient.XRevRangeN(stream, "+", "-", 1).Result()
	if err != nil {
		return nil, err
	}

	for _, row := range res {
		str, ok := row.Values["payload"].(string)
		if !ok {
			glog.Errorln(err)
			n.RedisClient.XDel(stream, row.ID)
			continue
		}
		item, err := DecodeForStream([]byte(str))
		if err != nil {
			glog.Errorln(err)
			n.RedisClient.XDel(stream, row.ID)
			continue
		}
		item.Id = row.ID
		return item, nil
	}
	return nil, errors.New("频道没有任何消息")
}
