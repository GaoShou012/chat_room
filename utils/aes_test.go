package utils

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"testing"
)

func TestAesEncrypt(t *testing.T) {
	key := []byte("NfmYCtxknKTN5wzZ")
	origData := make(map[string]interface{})
	origData["username"] = "abc"
	origData["nickname"] = "高手"

	var crypted []byte
	{
		j, err := json.Marshal(origData)
		if err != nil {
			glog.Errorln(err)
			return
		}
		data, err := AesEncrypt(j, key, nil, AesModeCBCPk5)
		if err != nil {
			glog.Errorln(err)
			return
		}
		crypted = data
	}

	{
		m := make(map[string]interface{})
		orig, err := AesDecrypt(crypted, key, nil, AesModeCBCPk5)
		if err != nil {
			glog.Errorln(err)
			return
		}
		if err := json.Unmarshal(orig, &m); err != nil {
			glog.Errorln(err)
			return
		}
		fmt.Println("orig data:", m)
	}
}
