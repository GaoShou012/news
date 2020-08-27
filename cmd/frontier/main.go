package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	uuid "github.com/satori/go.uuid"
	"im/config"
	proto_room "im/proto/room"
	"im/src/frontier"
	"im/src/im"
	"im/src/news"
	"im/utils"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/*
	监听kafka,group是根据frontierId进行监听
	所以每个frontierId都必须要不一样
*/

var loadFlags = micro.Action(func(c *cli.Context) error {
	return nil
})

type Handler struct {
	proto_room.RoomService
}

func (h *Handler) OnOpen(conn frontier.Conn) error {
	// 进行url验证连接
	// 并且获取acl列表
	return nil
}

func (h *Handler) OnMessage(conn frontier.Conn, message []byte) {
	msg := &im.Message{}
	err := json.Unmarshal(message, msg)
	if err != nil {
		return
	}
	switch msg.Head.BusinessType {
	case "news":
		news.News.OnMessage(conn, msg)
		break
	case "chatRoom":
		break
	}
}

func (h *Handler) OnClose(conn frontier.Conn) {

}

func main() {
	service := micro.NewService(
		micro.Flags(),
	)
	service.Init(loadFlags)

	utils.Micro.Init(service)
	utils.Micro.LoadSource()
	utils.Micro.LoadConfigMust(config.FrontierConfig)

	id := uuid.NewV4().String()
	addr := service.Server().Options().Address
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	protocol := &frontier.ProtocolWs{
		WriterBufferSize:     1024,
		ReaderBufferSize:     1024,
		WriterMessageTimeout: 20 * time.Millisecond,
		ReaderMessageTimeout: 20 * time.Millisecond,
	}
	handler := &Handler{RoomService: config.RoomServiceConfig.ServiceClient()}
	fr := &frontier.Frontier{
		Id:               id,
		Debug:            config.FrontierConfig.Debug,
		HeartbeatTimeout: config.FrontierConfig.HeartbeatTimeout,
		Protocol:         protocol,
		Handler:          handler,
	}
	if err := fr.Init(ctx, addr); err != nil {
		panic(err)
	}
	if err := fr.Start(); err != nil {
		panic(err)
	}
	go frontierConfigWatcher(fr)

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
		if err := fr.Stop(); err != nil {
			glog.Errorf("stop server error: %v", err)
		}
		break
	}
}

func frontierConfigWatcher(fr *frontier.Frontier) {
	watcher, err := utils.Micro.Config.Watch("micro", "config", "frontier")
	if err != nil {
		panic(err)
	}
	for {
		val, err := watcher.Next()
		if err != nil {
			panic(err)
		}
		conf := val.Bytes()
		fmt.Println("To update frontier config", string(conf))
		err = json.Unmarshal(conf, config.FrontierConfig)
		if err != nil {
			panic(err)
		}
		fr.Debug = config.FrontierConfig.Debug
		fr.HeartbeatTimeout = config.FrontierConfig.HeartbeatTimeout
	}
}
