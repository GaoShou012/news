package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
	"wchatv1/config"
	proto_room "wchatv1/proto/room"
	"wchatv1/utils"
)

func main() {
	utils.Micro.InitV2()
	utils.Micro.LoadSource()

	utils.Micro.LoadConfigMust(config.RoomServiceConfig)
	fmt.Println("room-service", config.RoomServiceConfig)

	passToken := &proto_room.PassToken{
		TenantCode: "9988",
		RoomCode:   "321",
		UserType:   "manager",
		UserId:     1,
		UserName:   "abc",
		UserThumb:  "123",
		UserTags:   "",
	}
	message := &proto_room.Message{
		Type:               "message.broadcast",
		Content:            "hello",
		ContentType:        "text",
		ClientMsgId:        123,
		ServerMsgId:        "",
		ServerMsgTimestamp: 0,
		From: &proto_room.From{
			Path:      "",
			PassToken: passToken,
		},
		To:       nil,
		RuledOut: nil,
	}

	wg := sync.WaitGroup{}
	wg.Add(1000)
	ch := make(chan *proto_room.Message, 100000)
	for i := 0; i < runtime.NumCPU()*30; i++ {
		go func() {
			for {
				message := <-ch
				req := &proto_room.BroadcastReq{Message: message}
				_, err := config.RoomServiceConfig.ServiceClient().Broadcast(context.TODO(), req)
				if err != nil {
					panic(err)
				}
				wg.Done()
			}
		}()
	}

	time.Sleep(time.Second * 3)
	fmt.Println(time.Now().Unix(), "begin")
	for i := 0; i < 1000; i++ {
		ch <- message
	}
	fmt.Println(time.Now().Unix(), "waiting")
	wg.Wait()
	fmt.Println(time.Now().Unix(), "end")
}
