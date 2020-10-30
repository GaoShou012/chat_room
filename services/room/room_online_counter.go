package room

import (
	"encoding/json"
	"fmt"
	"time"
)

type CountData struct {
	Count uint64
	Time  time.Time
}

func SetRoomUsersCount(frontierId uint64, tenantCode string, roomCode string, count uint64) error {
	fid := fmt.Sprintf("%d", frontierId)
	roomKey := fmt.Sprintf("room:online:count:%s:%s", tenantCode, roomCode)
	data := &CountData{
		Count: count,
		Time:  time.Now(),
	}
	j,err := json.Marshal(data)
	if err != nil {
		return err
	}
	RedisClusterClient.HSet(roomKey, fid, j)
	return nil
}
func GetRoomUsersCount(tenantCode string, roomCode string) (count uint64, err error) {
	roomKey := fmt.Sprintf("room:online:count:%s:%s", tenantCode, roomCode)
	res, err := RedisClusterClient.HGetAll(roomKey).Result()
	if err != nil {
		return 0,err
	}
	data := &CountData{}
	now := time.Now()
	for _, v := range res {
		err := json.Unmarshal([]byte(v), data)
		if err != nil {
			return 0,err
		}
		if now.Sub(data.Time) < 60 {
			count += data.Count
		}
	}
	return count,err
}
