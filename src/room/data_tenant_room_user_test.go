package room

import (
	"fmt"
	"github.com/golang/glog"
	"testing"
)

func TestTenantRoomUser_Set(t *testing.T) {
	fmt.Println("测试租户-房间-用户 set")
	tenantCode := "we_testing"
	roomCode := "room1"
	userId := uint64(100)
	uuid := "123123"

	// 循环设置40个用户
	tru := newTenantRoomUser()

	for i := 0; i < 40; i++ {
		err := tru.Set(tenantCode, roomCode, userId+uint64(i), uuid)
		if err != nil {
			glog.Errorln(err)
			return
		}
	}

	{
		err := tru.Set(tenantCode, "room2", userId, uuid)
		if err != nil {
			glog.Errorln(err)
			return
		}
	}
}

func TestTenantRoomUser_Del(t *testing.T) {
	fmt.Println("测试租户-房间-用户 del")

	tenantCode := "we_testing"
	roomCode := "room2"
	userId := uint64(100)
	uuid := "123123a"

	tru := newTenantRoomUser()

	err := tru.Del(tenantCode, roomCode, userId, uuid)
	if err != nil {
		glog.Errorln(err)
		return
	}
}

func TestTenantRoomUser_Range(t *testing.T) {
	fmt.Println("测试租户-房间-用户 range")
	tenantCode := "we_testing"
	roomCode := "room1"

	tru := newTenantRoomUser()
	res, err := tru.Range(tenantCode, roomCode, 0, 20)
	if err != nil {
		glog.Errorln(err)
		return
	}
	for _, row := range res {
		fmt.Println(row)
	}
}
