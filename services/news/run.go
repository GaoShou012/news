package news

import (
	"github.com/go-redis/redis/v7"
	"github.com/golang/glog"
	proto_news "gitlab.cptprod.net/wchat-backend/im-lib/proto/news"
	"im/config"
	"im/utils"
	"os"
	"runtime"
)

func RunService() {
	utils.Micro.LoadConfigMust(config.RedisClusterConfig)
	utils.RedisClusterClient = initRedis(config.RedisClusterConfig.Addr, config.RedisClusterConfig.Password)

	service := &Service{}
	if err := proto_news.RegisterNewsServiceHandler(utils.Micro.Service.Server(), service); err != nil {
		glog.Errorln(err)
		os.Exit(1)
	}
	if err := utils.Micro.Service.Run(); err != nil {
		glog.Errorln(err)
		os.Exit(1)
	}
}

func initRedis(addr []string, password string) *redis.ClusterClient {
	minSize := runtime.NumCPU() * 10
	poolSize := minSize
	minIdleConns := minSize

	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:              addr,
		MaxRedirects:       0,
		ReadOnly:           false,
		RouteByLatency:     false,
		RouteRandomly:      false,
		ClusterSlots:       nil,
		OnNewNode:          nil,
		Dialer:             nil,
		OnConnect:          nil,
		Username:           "",
		Password:           password,
		MaxRetries:         0,
		MinRetryBackoff:    0,
		MaxRetryBackoff:    0,
		DialTimeout:        0,
		ReadTimeout:        0,
		WriteTimeout:       0,
		NewClient:          nil,
		PoolSize:           poolSize,
		MinIdleConns:       minIdleConns,
		MaxConnAge:         0,
		PoolTimeout:        0,
		IdleTimeout:        0,
		IdleCheckFrequency: 0,
		TLSConfig:          nil,
	})
}
