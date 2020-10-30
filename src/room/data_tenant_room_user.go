package room

import (
	"errors"
	"github.com/go-redis/redis/v7"
)

type TenantRoomUser struct {
	RedisClient             *redis.ClusterClient
	TenantRoomUserMap       *TenantRoomUserMap
	TenantRoomUserSortedSet *TenantRoomUserSortedSet
}

func (t *TenantRoomUser) Set(tenantCode string, roomCode string, userId uint64, uuid string) error {
	{
		err := t.TenantRoomUserMap.Set(tenantCode, roomCode, userId, uuid)
		if err != nil {
			return err
		}
	}

	if err := t.TenantRoomUserSortedSet.Set(tenantCode, roomCode, userId); err != nil {
		return err
	}

	return nil
}

func (t *TenantRoomUser) Del(tenantCode string, roomCode string, userId uint64, uuid string) error {
	{
		id, err := t.TenantRoomUserMap.Get(tenantCode, roomCode, userId)
		if err != nil {
			return err
		}
		if id != uuid {
			return errors.New("id is not equal uuid")
		}
	}

	{
		err := t.TenantRoomUserMap.Del(tenantCode, roomCode, userId)
		if err != nil {
			return err
		}
	}

	{
		if err := t.TenantRoomUserSortedSet.Del(tenantCode, roomCode, userId); err != nil {
			return err
		}
	}

	return nil
}

func (t *TenantRoomUser) Count(tenantCode string, roomCode string) (int64, error) {
	return t.TenantRoomUserSortedSet.Count(tenantCode, roomCode)
}

func (t *TenantRoomUser) Range(tenantCode string, roomCode string, page uint64, pageSize uint64) ([]*TenantUserInfo, error) {
	// get users' id from sorted set
	idList, err := t.TenantRoomUserSortedSet.Range(tenantCode, roomCode, page, pageSize)
	if err != nil {
		return nil, err
	}

	var users []*TenantUserInfo

	// get users' info
	{
		for _, userId := range idList {
			user := &TenantUserInfo{}
			if err := user.Get(tenantCode, userId, t.RedisClient); err != nil {
				return nil, err
			}
			users = append(users, user)
		}
	}

	return users, nil
}
