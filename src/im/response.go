package im

import (
	"encoding/json"
	"github.com/golang/glog"
	"im/src/frontier"
)

func ResponseError(conn frontier.Conn, message *Message, err error) {
	j, err := json.Marshal(message)
	if err != nil {
		glog.Errorln(err)
		return
	}
	conn.Sender(j)
}

func Ack(conn frontier.Conn, message *Message) {
	j, err := json.Marshal(message)
	if err != nil {
		glog.Errorln(err)
		return
	}
	conn.Sender(j)
}
