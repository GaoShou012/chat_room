package tenant_api

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"testing"
	"wchatv1/utils"
)

func TestAuth_Verify(t *testing.T) {
	// 生成token

	var origData []byte
	{
		operator := &Operator{TenantCode: "bob"}
		j, err := json.Marshal(operator)
		if err != nil {
			glog.Errorln(err)
			return
		}
		origData = j
	}

	token, err := utils.AesEncrypt(origData, []byte(TokenKey), nil, utils.AesModeCBCPk5)
	if err != nil {
		glog.Errorln(err)
		return
	}
	fmt.Println("token:", token)
}
