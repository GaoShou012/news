package news

import (
	"fmt"
	"testing"
)

var testingVars struct {
	handler *Handler
	cli1    *Client
	cli2    *Client
}

func TestHandler_Subscribe(t *testing.T) {
	fmt.Println("测试订阅")

	h := &Handler{}
	h.OnInit()
	testingVars.handler = h

	{
		cli := NewClient(nil)
		subscribe := []string{"ch1", "ch2"}
		h.subscribe(cli, subscribe)
		testingVars.cli1 = cli
	}

	{
		cli := NewClient(nil)
		subscribe := []string{"ch2", "ch3"}
		h.subscribe(cli, subscribe)
		testingVars.cli2 = cli
	}

	for key, val := range h.channels {
		fmt.Printf("频道 %s 人数 :%d\n", key, val.Len())
	}
}

func TestHandler_Unsubscribe(t *testing.T) {
	fmt.Println("测试取消订阅")

	handler := testingVars.handler

	{
		fmt.Println("取消前")
		for key, val := range handler.channels {
			fmt.Printf("频道 %s 人数 :%d\n", key, val.Len())
		}
	}

	{
		fmt.Println("cli1")
		cli := testingVars.cli1
		handler.unsubscribe(cli)
		for key, val := range handler.channels {
			fmt.Printf("频道 %s 人数 :%d\n", key, val.Len())
		}
	}

	{
		fmt.Println("cli2")
		cli := testingVars.cli2
		handler.unsubscribe(cli)
		for key, val := range handler.channels {
			fmt.Printf("频道 %s 人数 :%d\n", key, val.Len())
		}
	}
}
