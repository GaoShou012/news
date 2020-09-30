package news

import (
	"context"
	"github.com/golang/glog"
	proto_news "gitlab.cptprod.net/wchat-backend/im-lib/proto/news"
	"im/utils"
	"time"
)

var _ proto_news.NewsServiceHandler = &Service{}

type Service struct {
}

func (s Service) GetSubList(ctx context.Context, req *proto_news.GetSubListReq, rsp *proto_news.GetSubListRsp) error {
	key := KeyOfSubListOfChannel(req.ChannelName)
	res, err := utils.RedisClusterClient.HGetAll(key).Result()
	if err != nil {
		glog.Errorln(err)
		return nil
	}
	for channelName, _ := range res {
		rsp.Subscribers = append(rsp.Subscribers, channelName)
	}
	return nil
}

func (s Service) Sub(ctx context.Context, req *proto_news.SubReq, rsp *proto_news.SubRsp) error {
	now := time.Now()
	pipe := utils.RedisClusterClient.TxPipeline()

	for _, channelName := range req.Channels {
		key := KeyOfSubListOfChannel(channelName)
		pipe.HSet(key, req.FrontierId, now.Unix())
	}
	_, err := pipe.Exec()
	if err != nil {
		glog.Errorln(err)
		return nil
	}

	return nil
}

func (s Service) Cancel(ctx context.Context, req *proto_news.CancelReq, rsp *proto_news.CancelRsp) error {
	pipe := utils.RedisClusterClient.TxPipeline()
	for _, channelName := range req.Channels {
		key := KeyOfSubListOfChannel(channelName)
		pipe.HDel(key, req.FrontierId)
	}
	_, err := pipe.Exec()
	if err != nil {
		glog.Errorln(err)
		return nil
	}

	return nil
}

func (s Service) GetNews(ctx context.Context, req *proto_news.GetNewsReq, rsp *proto_news.GetNewsRsp) error {
	panic("implement me")
}
