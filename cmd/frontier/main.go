package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wchatv1/config"
	proto_room "wchatv1/proto/room"
	"wchatv1/src/frontier"
	"wchatv1/src/ider"
	"wchatv1/src/im"
	"wchatv1/src/news"
	"wchatv1/utils"
)

/*
	监听kafka,group是根据frontierId进行监听
	所以每个frontierId都必须要不一样
*/

var (
	frontierId int
)

var loadFlags = micro.Action(func(c *cli.Context) error {
	frontierId = c.Int("frontier_id")
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
		micro.Name(config.FrontierServiceConfig.ServiceName()),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Flags(
			&cli.IntFlag{Name: "frontier_id", Usage: "The frontier ID"},
		),
	)
	service.Init(loadFlags)

	utils.Micro.InitV2()
	utils.Micro.LoadSource()

	fmt.Printf("Frontier ID = %d\n", frontierId)
	// 初始化雪花算法
	ider.InitSnowFlake(frontierId)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addr := service.Options().Server.Options().Address
	ln, err := (&net.ListenConfig{KeepAlive: time.Second * 60}).Listen(ctx, "tcp", addr)
	if err != nil {
		panic(err)
	}

	protocol := &frontier.ProtocolWs{
		WriterBufferSize:     1024,
		ReaderBufferSize:     1024,
		WriterMessageTimeout: 20 * time.Millisecond,
		ReaderMessageTimeout: 20 * time.Millisecond,
	}
	handler := &Handler{RoomService: config.RoomServiceConfig.ServiceClient()}
	fr := &frontier.Frontier{
		Id:               frontierId,
		Ln:               ln,
		Heartbeat:        true,
		HeartbeatTimeout: 90,
		Protocol:         protocol,
		Handler:          handler}
	err = fr.Init()
	if err != nil {
		panic(err)
	}
	err = fr.Start()
	if err != nil {
		panic(err)
	}

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
