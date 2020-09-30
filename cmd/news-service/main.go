package main

import (
	"github.com/golang/glog"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	proto_news "gitlab.cptprod.net/wchat-backend/im-lib/proto/news"
	"im/services/news"
	"im/utils"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	prometheusAddr string
)

var loadFlags = micro.Action(func(c *cli.Context) error {
	prometheusAddr = c.String("prometheus_address")
	return nil
})

func main() {
	service := micro.NewService(
		micro.Name(proto_news.ServiceName),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Flags(
			&cli.StringFlag{Name: "prometheus_address", Usage: "The prometheus service"},
		),
	)
	service.Init(loadFlags)

	utils.Micro.Init(service)
	utils.Micro.LoadSource()

	go utils.Prometheus(prometheusAddr)
	go news.RunService()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		switch s := <-c; s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			glog.Infof("got signal %s; stop server", s)
		case syscall.SIGHUP:
			glog.Infof("got signal %s; go to deamon", s)
			continue
		}
		break
	}
}
