package config

import (
	"sync"
	proto_room "wchatv1/proto/room"
	"wchatv1/utils"
)

var (
	_           utils.MicroServiceConfig = &businessAcl{}
	BusinessAcl                          = &businessAcl{}
)

type businessAcl struct {
	Topic string `json:"topic"`

	serviceClientInit sync.Once
	serviceClient     proto_room.RoomService
}

func (c *businessAcl) ConfigKey() string {
	return "room-service"
}

func (c *businessAcl) ServiceName() string {
	return "micro.service.room-service"
}

func (c *businessAcl) ServiceClient() proto_room.RoomService {
	c.serviceClientInit.Do(func() {
		c.serviceClient = proto_room.NewRoomService(c.ServiceName(), utils.Micro.Service.Client())
	})
	return c.serviceClient
}
