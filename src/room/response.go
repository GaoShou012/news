package room

import (
	"fmt"
	proto_room "wchatv1/proto/room"
	"wchatv1/src/frontier"
)

func RspUsersCountInRoom(conn frontier.Conn, count int) {
	//msg := proto_room.Message{
	//	Type:        "",
	//	Content:     "",
	//	ContentType: "",
	//	ClientMsgId: 0,
	//	ServerMsgId: "",
	//	From:        nil,
	//	To:          nil,
	//	RuledOut:    nil,
	//}
}

func RspMessageOk(conn frontier.Conn, msgType string, cliMsgId uint64, srvMsgId string) error {
	msg := &proto_room.Message{
		Type:               fmt.Sprintf("%s.state.ok", msgType),
		Content:            "",
		ContentType:        "",
		ClientMsgId:        cliMsgId,
		ServerMsgId:        srvMsgId,
		ServerMsgTimestamp: 0,
		From:               nil,
		To:                 nil,
		RuledOut:           nil,
	}
	bts, err := Codec.Encode(msg)
	if err != nil {
		return err
	}

	return conn.Writer(bts)
}
func RspMessageNok(conn frontier.Conn, msgType string, cliMsgId uint64, srvMsgId string, content string) error {
	msg := &proto_room.Message{
		Type:        fmt.Sprintf("%s.state.nok", msgType),
		Content:     content,
		ContentType: ContentTypeText,
		ClientMsgId: cliMsgId,
		ServerMsgId: srvMsgId,
		From:        nil,
		To:          nil,
		RuledOut:    nil,
	}
	bts, err := Codec.Encode(msg)
	if err != nil {
		return err
	}

	return conn.Writer(bts)
}

func RspSendMessageOk(conn frontier.Conn, cliMsgId uint64, srvMsgId string) error {
	msg := proto_room.Message{
		Type:        TypeMessageStateOk,
		Content:     "",
		ContentType: ContentTypeText,
		ClientMsgId: cliMsgId,
		ServerMsgId: srvMsgId,
		From:        nil,
		To:          nil,
		RuledOut:    nil,
	}
	bts, err := Codec.Encode(&msg)
	if err != nil {
		return err
	}

	return conn.Writer(bts)
}
func RspSendMessageNok(conn frontier.Conn, cliMsgId uint64, srvMsgId string, content string) error {
	msg := proto_room.Message{
		Type:        TypeMessageStateNok,
		Content:     content,
		ContentType: ContentTypeText,
		ClientMsgId: cliMsgId,
		ServerMsgId: srvMsgId,
		From:        nil,
		To:          nil,
		RuledOut:    nil,
	}
	bts, err := Codec.Encode(&msg)
	if err != nil {
		return err
	}

	return conn.Writer(bts)
}
