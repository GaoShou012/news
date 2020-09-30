package news

import (
	"im/utils"
	"sync"
)

const (
	ServiceName = "micro.service.news-service"
)

var (
	serviceClientInit sync.Once
	serviceClient     NewsService
)

func ServiceClient() NewsService {
	serviceClientInit.Do(func() {
		serviceClient = NewNewsService(ServiceName, utils.Micro.Service.Client())
	})
	return serviceClient
}
