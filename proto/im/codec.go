package im

import (
	"github.com/golang/glog"
	"google.golang.org/protobuf/proto"
)

func Encode(message *Message) []byte {
	//data, err := json.Marshal(message)
	data, err := proto.Marshal(message)
	if err != nil {
		glog.Errorln(message)
		panic(err)
	}
	return data
}

func Decode(data []byte) (message *Message, err error) {
	message = &Message{}
	err = proto.Unmarshal(data, message)
	return
}
