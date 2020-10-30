package utils

import (
	"fmt"
	"github.com/golang/glog"
	"testing"
	"time"
)

func TestDisSortedSet_Add(t *testing.T) {
	fmt.Println("测试DisSortedSet Add")

	key := fmt.Sprintf("testing:dis_sorted_set")
	disSortedSet := NewDisSortedSet(newRedis())
	if err := disSortedSet.Add(key, "abc"); err != nil {
		glog.Errorln(err)
		return
	}
	time.Sleep(time.Second*5)
	if err := disSortedSet.Add(key, "abc1"); err != nil {
		glog.Errorln(err)
		return
	}
}

func TestDisSortedSet_GetValid(t *testing.T) {
	fmt.Println("测试DisSortedSet GetValid")

	key := fmt.Sprintf("testing:dis_sorted_set")
	disSortedSet := NewDisSortedSet(newRedis())
	res, err := disSortedSet.GetValid(key, 3)
	if err != nil {
		glog.Errorln(err)
		return
	}
	for _, row := range res {
		fmt.Println(row)
	}
}

func TestDisSortedSet_Count(t *testing.T) {
	fmt.Println("测试DisSortedSet Count")

	key := fmt.Sprintf("testing:dis_sorted_set")
	disSortedSet := NewDisSortedSet(newRedis())
	res, err := disSortedSet.Count(key, 90)
	if err != nil {
		glog.Errorln(err)
		return
	}
	fmt.Println("count:", res)
}
