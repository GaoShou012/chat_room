package room

import (
	"fmt"
	"github.com/golang/glog"
	"strconv"
	"strings"
	"sync"
)

func ruledOutSender(clientsIndex *sync.Map, fromPath string, ruledOut map[*Client]bool) {
	if ruledOut == nil {
		return
	}

	// 拆分字符串
	arr := strings.Split(fromPath, ".")
	if len(arr) != 5 {
		//glog.Errorln("is not match 5 fields")
		return
	}

	// fid,connId,tenantCode,roomCode,userId
	_, cid, tenantCode, _, uid := arr[0], arr[1], arr[2], arr[3], arr[4]

	connId, err := strconv.Atoi(cid)
	if err != nil {
		glog.Errorln("ConnId 类型转换失败")
		return
	}
	userId, err := strconv.Atoi(uid)
	if err != nil {
		glog.Errorln("UserId 类型转换失败")
		return
	}

	// 客户端索引值，使用connId，导出连接
	val, ok := clientsIndex.Load(connId)
	if !ok {
		fmt.Println("load failed")
		return
	}
	cli, ok := val.(*Client)
	if !ok {
		fmt.Println("duanyan is failed")
		return
	}

	// 判断连接是否同一个用户
	ctx := cli.conn.GetContext().(*ConnContext)
	if ctx.TenantCode != tenantCode {
		fmt.Println("tenantCode is not match")
		//return
	}
	if ctx.UserId != uint64(userId) {
		fmt.Println("user id is not match")
		//return
	}

	ruledOut[cli] = true
}
