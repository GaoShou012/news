package room

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/golang/glog"
	"sync"
	proto_room "wchatv1/proto/room"
)

var _ proto_room.RoomServiceHandler = &Service{}

type Service struct {
	Key                 []byte
	BroadcastToFrontier chan *proto_room.Message

	mutex sync.Mutex
}

func (s *Service) Init() {

}

func (s *Service) MakePassToken(ctx context.Context, req *proto_room.MakePassTokenReq, rsp *proto_room.MakePassTokenRsp) error {
	passToken := &proto_room.PassToken{
		TenantCode: req.TenantCode,
		RoomCode:   req.RoomCode,
		UserType:   req.UserType,
		UserId:     req.UserId,
		UserName:   req.UserName,
		UserThumb:  req.UserThumb,
		UserTags:   req.UserTags,
	}

	str, err := Codec.EncodePassToken(s.Key, passToken)
	if err != nil {
		return err
	}

	rsp.Token = str
	return nil
}

func (s *Service) ViewPassToken(ctx context.Context, req *proto_room.ViewPassTokenReq, rsp *proto_room.ViewPassTokenRsp) error {
	passToken, err := Codec.DecodePassToken(s.Key, req.Token)
	if err != nil {
		return err
	}
	rsp.PassToken = passToken
	return nil
}

func (s *Service) Join(ctx context.Context, req *proto_room.JoinReq, rsp *proto_room.JoinRsp) error {
	passToken := req.PassToken

	// 增加在线人数数量
	// 发出加入通知
	tenantCode := passToken.TenantCode
	roomCode := passToken.RoomCode
	num, err := IncrUserCount(tenantCode, roomCode)
	if err != nil {
		fmt.Println(err)
	}
	s.BroadcastToFrontier <- NotificationUserJoin(tenantCode, roomCode, passToken)
	s.BroadcastToFrontier <- NotificationUsersCount(tenantCode, roomCode, num)
	return nil
}

func (s *Service) Leave(ctx context.Context, req *proto_room.LeaveReq, rsp *proto_room.LeaveRsp) error {
	passToken := req.PassToken
	tenantCode := passToken.TenantCode
	roomCode := passToken.RoomCode
	num, err := DecrUserCount(tenantCode, roomCode)
	if err != nil {
		fmt.Println(err)
	}
	s.BroadcastToFrontier <- NotificationUserLeave(tenantCode, roomCode, passToken)
	s.BroadcastToFrontier <- NotificationUsersCount(tenantCode, roomCode, num)

	return nil
}

func (s *Service) GetUsersCount(ctx context.Context, req *proto_room.GetUsersCountReq, rsp *proto_room.GetUsersCountRsp) error {
	count, err := GetRoomUsersCount(req.TenantCode, req.RoomCode)
	if err != nil {
		return err
	}

	rsp.Count = count
	return nil
}

func (s *Service) SetUsersCount(ctx context.Context, req *proto_room.SetUsersCountReq, rsp *proto_room.SetUsersCountRsp) error {
	err := SetRoomUsersCount(req.FrontierId, req.TenantCode, req.RoomCode, req.Count)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Broadcast(ctx context.Context, req *proto_room.BroadcastReq, rsp *proto_room.BroadcastRsp) error {
	if req.Message.From == nil {
		return fmt.Errorf("message.from is nil")
	}

	message := req.Message
	passToken := message.From.PassToken

	aclCtx, err := Acl(message)
	if err != nil {
		rsp.Message = MsgFormatStateNok(req.Message.Type, req.Message.ClientMsgId, err.Error())
		return nil
	}

	switch req.Message.Type {
	case "message.broadcast":
		message.To = &proto_room.To{Path: fmt.Sprintf("%s.%s", passToken.TenantCode, passToken.RoomCode)}
		mid, err := Stream(passToken.TenantCode, passToken.RoomCode, message)
		if err != nil {
			return err
		}
		message.ServerMsgId = mid
		s.BroadcastToFrontier <- message

		rsp.Message = MsgFormatStateOk(req.Message.Type, req.Message.ClientMsgId)
		rsp.Message.ServerMsgId = mid
		break
	case "message.cancel":
		targetMsg, ok := aclCtx.(*proto_room.Message)
		if !ok {
			return fmt.Errorf("aclCtx Assert payload is failed")
		}
		DelMessageById(passToken.TenantCode, passToken.RoomCode, targetMsg.ServerMsgId)

		msgId := targetMsg.ServerMsgId
		msg := MsgFormatMessageCancel(msgId, passToken)
		Stream(passToken.TenantCode, passToken.RoomCode, msg)
		s.BroadcastToFrontier <- msg

		rsp.Message = MsgFormatStateOk(req.Message.Type, req.Message.ClientMsgId)
		rsp.Message.ServerMsgId = msgId
		break
	case "acl.message.broadcast":
		permissionKey := "message.broadcast"
		permissionVal := message.Content

		if message.To.Path == "all" {
			if permissionVal != "0" && permissionVal != "1" {
				rsp.Message = MsgFormatStateNok(req.Message.Type, req.Message.ClientMsgId, "Content内容不符合要求 0 or 1")
				return nil
			}
			err := SetRoomAcl(
				passToken.TenantCode, passToken.RoomCode,
				permissionKey, permissionVal,
			)
			if err != nil {
				return err
			}
		} else {
			err := SetUserAcl(
				passToken.TenantCode, passToken.RoomCode, passToken.UserId,
				permissionKey, permissionVal,
			)
			if err != nil {
				return err
			}
		}

		rsp.Message = MsgFormatStateOk(req.Message.Type, req.Message.ClientMsgId)
		break
	default:
		rsp.Message = MsgFormatStateNok(req.Message.Type, req.Message.ClientMsgId, "未处理的操作")
	}

	//// TODO (修正广播消息上下文路由)
	//if req.Message.Type == "broadcast" {
	//	passToken := req.Message.From.PassToken
	//	switch req.Message.From.PassToken.UserType {
	//	case "user":
	//		req.Message.To.Path = fmt.Sprintf("room.%s.%s", passToken.TenantCode, passToken.RoomCode)
	//		break
	//	case "manager":
	//		req.Message.To.Path = fmt.Sprintf("room.%s.%s", passToken.TenantCode, passToken.RoomCode)
	//		break
	//	}
	//}
	//
	//// TODO （检查Redis，用户是否已经被禁言）
	//// TODO (鉴权、敏感词过滤)
	//if err := CheckAuthority(req.Message); err != nil {
	//	rsp.Code = 1
	//	rsp.Message = err.Error()
	//	return rsp, nil
	//}
	//if err := CheckWords(req.Message.Content); err != nil {
	//	rsp.Code = 1
	//	rsp.Message = err.Error()
	//	return rsp, nil
	//}
	return nil
}

func (s *Service) Record(ctx context.Context, req *proto_room.RecordReq, rsp *proto_room.RecordRsp) error {
	stream := fmt.Sprintf("im:stream:%s:%s", req.TenantCode, req.RoomCode)
	lastMessageId, count := req.LastMessageId, req.Count
	if lastMessageId == "" {
		lastMessageId = "0"
	}
	if count == 0 || count > 1000 {
		count = 20
	}

	res, err := RedisClusterClient.XRead(&redis.XReadArgs{
		Streams: []string{stream, lastMessageId},
		Count:   int64(count),
		Block:   -1,
	}).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}

	var record []*proto_room.Message
	for _, val := range res {
		for _, message := range val.Messages {
			payload, ok := message.Values["payload"].(string)
			if !ok {
				glog.Errorln("Assert payload is failed")
				RedisClusterClient.XDel(stream, message.ID)
				continue
			}

			msg, err := Codec.DecodeMessage([]byte(payload))
			if err != nil {
				glog.Errorln(err)
				RedisClusterClient.XDel(stream, message.ID)
				continue
			}

			if msg.Type == "message.broadcast" || msg.Type == "message.p2p" {
				msg.ServerMsgId = message.ID
			}

			record = append(record, msg)
		}
	}

	rsp.Record = record
	return nil
}
