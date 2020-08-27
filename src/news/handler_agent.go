package news

import (
	"fmt"
	"github.com/golang/glog"
	"im/src/frontier"
	"im/src/im"
	"runtime"
	"sync"
)

type handlerAgent struct {
	clients  sync.Map
	handlers []*handler
}

func (h *handlerAgent) init() {
	for i := 0; i < runtime.NumCPU(); i++ {
		handler := &handler{}
		handler.init()
		h.handlers = append(h.handlers, handler)
	}
}
func (h *handlerAgent) SendMessage(conn frontier.Conn, message *im.Message) {
	msg, err := im.Encode(message)
	if err != nil {
		glog.Errorln(err)
		return
	}
	conn.Sender(msg)
}

func (h *handlerAgent) AddConn(conn frontier.Conn) {
	id := conn.GetId()
	h.handlers[id%runtime.NumCPU()].add
}
func (h *handlerAgent) DelConn() {

}

func (h *handlerAgent) UploadSubscribeList(conn frontier.Conn, message *im.Message) {
	id := conn.GetId()
	val, ok := h.clients.Load(id)
	if !ok {
		message.Head.Status = 1
		message.Body = "Cli不存在"
		h.SendMessage(conn, message)
		return
	}
	cli := val.(*Client)
	subList, ok := message.Body.([]string)
	if !ok {
		message.Head.Status = 1
		message.Body = fmt.Sprintf("监听列表转换失败")
		h.SendMessage(conn, message)
		return
	}
	if len(subList) > 255 {
		message.Head.Status = 1
		message.Body = fmt.Sprintf("监听列表超出数量限制")
		h.SendMessage(conn, message)
		return
	}
	h.handlers[id%runtime.NumCPU()].UploadSubList(cli, subList)
}
