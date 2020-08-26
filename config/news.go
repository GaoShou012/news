package config

import (
	proto_news "im/proto/news"
	"im/utils"
	"sync"
)

var (
	_                 utils.MicroServiceConfig = &newsServiceConfig{}
	NewsServiceConfig                          = &newsServiceConfig{}
)

type newsServiceConfig struct {
	Topic string `json:"topic"`

	serviceClientInit sync.Once
	serviceClient     proto_news.NewsService
}

func (c *newsServiceConfig) ConfigKey() string {
	return "room-service"
}

func (c *newsServiceConfig) ServiceName() string {
	return "micro.service.room-service"
}

func (c *newsServiceConfig) ServiceClient() proto_news.NewsService {
	c.serviceClientInit.Do(func() {
		c.serviceClient = proto_news.NewNewsService(c.ServiceName(), utils.Micro.Service.Client())
	})
	return c.serviceClient
}
