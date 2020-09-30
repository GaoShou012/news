package main

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	proto_news "im/proto/news"
	"im/utils"
)

func main() {
	utils.Micro.Init(nil)
	utils.Micro.LoadSource()
	cli := proto_news.ServiceClient()

	{
		req := &proto_news.SubReq{
			FrontierId: "123555",
			Channels:   []string{"123", "4444"},
		}
		rsp, err := cli.Sub(context.TODO(), req)
		if err != nil {
			glog.Errorln(err)
			return
		}
		fmt.Println(rsp)
	}

	{
		req := &proto_news.GetSubListReq{ChannelName: "123"}
		rsp, err := cli.GetSubList(context.TODO(), req)
		if err != nil {
			glog.Errorln(err)
			return
		}
		for _, frontierId := range rsp.Subscribers {
			fmt.Println("f=", frontierId)
		}
	}

	{
		req := &proto_news.CancelReq{
			FrontierId: "123555",
			Channels:   []string{"123"},
		}
		rsp,err := cli.Cancel(context.TODO(),req)
		if err != nil {
			glog.Errorln(err)
			return
		}
		fmt.Println(rsp)
	}

	{
		req := &proto_news.GetSubListReq{ChannelName: "123"}
		rsp, err := cli.GetSubList(context.TODO(), req)
		if err != nil {
			glog.Errorln(err)
			return
		}
		for _, frontierId := range rsp.Subscribers {
			fmt.Println("f=", frontierId)
		}
	}
}
