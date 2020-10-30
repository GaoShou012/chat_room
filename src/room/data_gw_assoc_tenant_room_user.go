package room

import (
	"errors"
	"github.com/golang/glog"
)

type GwAssocTenantRoomsUsers struct {
	GwAssocTenantRoomMap       *GwAssocTenantRoomMap
	GwAssocTenantRoomsUsersMap *GwAssocTenantRoomUserMap
}

/*
	设置的时候，不需要检查，之前保存的uuid
	设置即生效
*/
func (g *GwAssocTenantRoomsUsers) Set(gwId string, tenantCode string, roomCode string, userId uint64, uuid string) error {
	{
		err := g.GwAssocTenantRoomMap.Set(gwId, tenantCode, roomCode)
		if err != nil {
			return err
		}
	}

	{
		err := g.GwAssocTenantRoomsUsersMap.Set(gwId, tenantCode, roomCode, userId, uuid)
		if err != nil {
			return err
		}
	}

	return nil
}

/*
	移除的时候，需要检查之前的uuid是否一致
	有可能在时序上出现混乱，所以必须要检查
*/
func (g *GwAssocTenantRoomsUsers) Del(gwId string, tenantCode string, roomCode string, userId uint64, uuid string) error {
	id, err := g.GwAssocTenantRoomsUsersMap.Get(gwId, tenantCode, roomCode, userId)
	if err != nil {
		return err
	}
	if id != uuid {
		return errors.New("uuid is different")
	}

	// 移除网关绑定的用户
	if err := g.GwAssocTenantRoomsUsersMap.Del(gwId, tenantCode, roomCode, userId); err != nil {
		return err
	}

	// 房间数量为0，移除网关绑定的房间数据
	num, err := g.GwAssocTenantRoomsUsersMap.Len(gwId, tenantCode, roomCode)
	if err != nil {
		return err
	}
	if num == 0 {
		if err := g.GwAssocTenantRoomsUsersMap.DelAll(gwId, tenantCode, roomCode); err != nil {
			glog.Errorln(err)
		}
		if err := g.GwAssocTenantRoomMap.DelAll(gwId); err != nil {
			glog.Errorln(err)
		}
	}

	return nil
}
