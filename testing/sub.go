package main

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	proto_im "im/proto/im"
	proto_news "im/proto/news"
	"time"
)

func main() {
	c, _, err := websocket.DefaultDialer.Dial("ws://192.168.56.101:2330?mid=1598252141521-0&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb29tQ29kZSI6Ijg4OWQiLCJ0ZW5hbnRDb2RlIjoieGdjcCIsInVzZXJJZCI6MzIxLCJ1c2VyTmFtZSI6ImFiY2NjIiwidXNlclRodW1iIjoiZGRkIiwidXNlclR5cGUiOiJtYW5hZ2VyIn0.QNSZ3_9sn1VK4ZANAMRoOQMvVnCZLfb2_quF9-dcO7E", nil)
	if err != nil {
		glog.Errorln(err)
		return
	}
	defer c.Close()

	go func() {
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				glog.Errorln(err)
				break
			}
			imMessage, err := proto_im.Decode(data)
			if err != nil {
				glog.Errorln(err)
				continue
			}
			fmt.Println(imMessage)
		}
	}()

	sub := &proto_news.Subscribe{Channels: []string{"channela","channelb"}}
	data := proto_news.Request(sub)
	c.WriteMessage(websocket.BinaryMessage, data)
	time.Sleep(time.Second * 3)
}
