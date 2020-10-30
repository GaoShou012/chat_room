package room

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	"testing"
	proto_room "wchatv1/proto/room"
)

func TestInitRedisClusterClient(t *testing.T) {
	addr := []string{"192.168.56.101:9001", "192.168.56.101:9002", "192.168.56.101:9003", "192.168.56.101:9004", "192.168.56.101:9005", "192.168.56.101:9006"}
	password := ""
	InitRedisClusterClient(addr, password)
}

func TestService_SetVirtualUserCount(t *testing.T) {
	service := &Service{}
	req := &proto_room.SetVirtualUserCountReq{
		TenantCode: "cc",
		RoomCode:   "cc3",
		Count:      100,
	}
	rsp := &proto_room.SetVirtualUserCountRsp{}
	err := service.SetVirtualUserCount(context.TODO(), req, rsp)
	if err != nil {
		glog.Errorln(err)
		return
	}
	fmt.Println(rsp.Code, rsp.Desc)
}

func TestService_GetVirtualUserCount(t *testing.T) {
	service := &Service{}
	req := &proto_room.GetVirtualUserCountReq{
		TenantCode: "cc",
		RoomCode:   "cc3",
	}
	rsp := &proto_room.GetVirtualUserCountRsp{}
	err := service.GetVirtualUserCount(context.TODO(), req, rsp)
	if err != nil {
		glog.Errorln(err)
		return
	}
	fmt.Println(rsp.Count)
}
