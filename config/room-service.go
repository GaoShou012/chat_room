package config

import (
	"sync"
	proto_room "wchatv1/proto/room"
	"wchatv1/utils"
)

var (
	_                 utils.MicroServiceConfig = &roomServiceConfig{}
	RoomServiceConfig                          = &roomServiceConfig{}
)

type roomServiceConfig struct {
	Topic string `json:"topic"`

	serviceClientInit sync.Once
	serviceClient     proto_room.RoomService
}

func (c *roomServiceConfig) ConfigKey() string {
	return "room-service"
}

func (c *roomServiceConfig) ServiceName() string {
	return "micro.service.room-service"
}

func (c *roomServiceConfig) ServiceClient() proto_room.RoomService {
	c.serviceClientInit.Do(func() {
		c.serviceClient = proto_room.NewRoomService(c.ServiceName(), utils.Micro.Service.Client())
	})
	return c.serviceClient
}
