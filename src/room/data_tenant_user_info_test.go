package room

import (
	"fmt"
	"github.com/golang/glog"
	"testing"
)

func TestTenantUserInfo_Set(t *testing.T) {
	fmt.Println("测试租户信息Set")

	info := TenantUserInfo{
		TenantCode: "we_testing",
		UserType:   "user",
		UserId:     100,
		Username:   "abc",
		Nickname:   "高手",
		UserThumb:  "",
		UserTags:   "",
	}
	err := info.Set(info.TenantCode, info.UserId, newRedis())
	if err != nil {
		glog.Errorln(err)
		return
	}

	// 循环生成40个用户信息
	for i := 0; i < 40; i++ {
		info := &TenantUserInfo{
			TenantCode: "we_testing",
			UserType:   "user",
			UserId:     100 + uint64(i),
			Username:   "abc",
			Nickname:   fmt.Sprintf("高手-%d", i),
			UserThumb:  "",
			UserTags:   "",
		}
		err := info.Set(info.TenantCode, info.UserId, newRedis())
		if err != nil {
			glog.Errorln(err)
			return
		}
	}
}

func TestTenantUserInfo_Get(t *testing.T) {
	fmt.Println("测试租户信息Get")

	info := TenantUserInfo{}
	err := info.Get("we_testing", 100, newRedis())
	if err != nil {
		glog.Errorln(err)
		return
	}
	fmt.Println("info", info)
}
