package news

import (
	"fmt"
	"testing"
)

func TestChannel_Add(t *testing.T) {
	fmt.Println("测试channel添加")

	ch := Channel{}
	cli := &Client{
		Conn:            nil,
		News:            nil,
		ChannelsAnchors: nil,
	}
	anchor := ch.Add(cli)
	fmt.Println("anchor", anchor)
	fmt.Println("ch clients count", ch.Clients.Len())
}

func TestChannel_Del(t *testing.T) {
	fmt.Println("测试channel删除")

	ch := Channel{}
	cli := &Client{}
	anchor := ch.Add(cli)
	fmt.Println("clients count before del", ch.Len())

	ch.Del(anchor)
	fmt.Println("clients count after del", ch.Len())
}
