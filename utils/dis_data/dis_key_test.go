package dis_data

import (
	"fmt"
	"github.com/golang/glog"
	"testing"
	"time"
)

func TestDisHash_Add(t *testing.T) {
	fmt.Println("测试dis_hash add")
	key := "testing:dis_hash"
	hash := newDisKey()
	err := hash.Add(key, "abc")
	if err != nil {
		glog.Errorln(err)
		return
	}
	time.Sleep(time.Second * 5)
}

func TestDisHash_GetValid(t *testing.T) {
	fmt.Println("测试dis_hash get valid")
	key := "testing:dis_hash"
	hash := newDisKey()
	{
		rows, err := hash.GetValid(key, 3)
		if err != nil {
			glog.Errorln(err)
			return
		}
		for _, row := range rows {
			fmt.Println(row)
		}
	}
}
