package room

import (
	"container/list"
	"context"
	"fmt"
	"github.com/golang/glog"
	"github.com/micro/go-micro/v2/client"
	"os"
	"sync"
	"time"
	"wchatv1/config"
	proto_room "wchatv1/proto/room"
)

type Room struct {
	isOpen     bool
	mutex      sync.Mutex
	wg         *sync.WaitGroup
	tenantCode string
	roomCode   string

	// 结束事件
	handlerDone chan bool
	// 消息广播
	broadcast chan *proto_room.Message
	// 房间事件
	event chan *Event

	// 客户端
	clientsIndex  *sync.Map
	clients       *list.List
	clientsAnchor map[*Client]*list.Element
}

func (r *Room) Init(tenantCode string, roomCode string, eventCacheSize int, broadcastCacheSize int) {
	r.isOpen = true
	r.wg = &sync.WaitGroup{}
	r.tenantCode, r.roomCode = tenantCode, roomCode
	r.broadcast = make(chan *proto_room.Message, broadcastCacheSize)
	r.event = make(chan *Event, eventCacheSize)
	r.handlerDone = make(chan bool)
	r.clientsIndex = &sync.Map{}
	r.clients = list.New()
	r.clientsAnchor = make(map[*Client]*list.Element)
	go func() {
		r.wg.Add(1)
		defer r.wg.Done()
		for {
			select {
			case <-r.handlerDone:
				return
			case message := <-r.broadcast:
				if message == nil {
					if r.isOpen == false {
						fmt.Println("room was closed")
					}
					fmt.Println("room message is nil")
					os.Exit(1)
				}
				Sender.Broadcast(r.clients, r.clientsIndex, message).Wait()
			case event := <-r.event:
				if event == nil {
					if r.isOpen == false {
						fmt.Println("room was closed")
					}
					fmt.Println("event is nil")
					os.Exit(1)
				}
				cli := event.Client
				if cli == nil {
					fmt.Println("cli is nil")
					os.Exit(1)
				}
				switch event.Type {
				case RoomEventJoin:
					if _, ok := r.clientsAnchor[cli]; ok {
						glog.Errorln("RoomEventJoin The client is exists already", cli)
						break
					}

					// 保存Client
					anchor := r.clients.PushBack(cli)
					r.clientsAnchor[cli] = anchor
					r.clientsIndex.Store(cli.conn.GetId(), cli)

					// 获取当前用户数量
					//ctx := cli.conn.GetContext().(*ConnContext)
					//serviceClient := services.NewRoomServiceClient()
					//req := &proto_room.GetUsersCountReq{TenantCode:ctx.TenantCode,RoomCode:ctx.RoomCode}
					//rsp,err := serviceClient.GetUsersCount(context.TODO(),req)
					//if err != nil {
					//	glog.Errorln(err)
					//	break
					//}

					// 投放到历史消息同步
					SyncRecord.NewClient(cli)
				case RoomEventLeave:
					anchor, ok := r.clientsAnchor[cli]
					if !ok {
						glog.Errorln("RoomEventLeave The client is not exists.", cli)
						break
					}

					cli.isClose = true
					r.clients.Remove(anchor)
					delete(r.clientsAnchor, cli)
					r.clientsIndex.Delete(cli.conn.GetId())

					ctx := cli.conn.GetContext().(*ConnContext)
					go func(ctx *ConnContext) {
						passToken := &proto_room.PassToken{
							TenantCode: ctx.TenantCode,
							RoomCode:   ctx.RoomCode,
							UserType:   ctx.UserType,
							UserId:     ctx.UserId,
							UserName:   ctx.UserName,
							UserThumb:  ctx.UserThumb,
							UserTags:   ctx.UserTags,
						}
						req := &proto_room.LeaveReq{
							PassToken: passToken,
						}

						var opss client.CallOption = func(o *client.CallOptions) {
							o.RequestTimeout = time.Second * 10
							o.DialTimeout = time.Second * 10
						}
						serviceClient := config.RoomServiceConfig.ServiceClient()
						_, err := serviceClient.Leave(context.TODO(), req, opss)
						if err != nil {
							glog.Errorln(err)
						}
					}(ctx)
				case RoomEventClientToDirectLine:
					if _, ok := r.clientsAnchor[cli]; !ok {
						glog.Errorln("RoomEventClientToDirectLine The client is not exists.", cli)
						break
					}

					for {
						msg := cli.GetCache()
						if msg == nil {
							break
						}

						if cli.CheckMessageContinuity(msg) {
							Sender.Point(cli, msg)
						}
					}
					cli.isCache = false

					//fmt.Println("to direct line", cli.conn.GetContext())
				}
			}
		}
	}()
}

func (r *Room) Close() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.handlerDone <- true
	r.wg.Wait()
	r.isOpen = false
	r.broadcast = nil
	r.event = nil
}

func (r *Room) ClientsCounter() int {
	return r.clients.Len()
}

func (r *Room) Join(cli *Client) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.isOpen == false {
		return
	}
	r.event <- &Event{
		Type:   RoomEventJoin,
		Client: cli,
	}
}
func (r *Room) Leave(cli *Client) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.isOpen == false {
		return
	}
	r.event <- &Event{
		Type:   RoomEventLeave,
		Client: cli,
	}
}
func (r *Room) ClientToDirectMode(cli *Client) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.isOpen == false {
		return
	}
	r.event <- &Event{
		Type:   RoomEventClientToDirectLine,
		Client: cli,
	}
}

func (r *Room) Broadcast(message *proto_room.Message) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.isOpen == false {
		return
	}
	r.broadcast <- message
}
