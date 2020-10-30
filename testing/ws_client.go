package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"sync"
	"time"
	proto_room "wchatv1/proto/room"
)

func main() {
	c, _, err := websocket.DefaultDialer.Dial("wss://im-frontier.weprod.net/?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb29tQ29kZSI6Ijg4OWQiLCJ0ZW5hbnRDb2RlIjoieHh4eCIsInVzZXJJZCI6NzI4MjA1NTQ3MjU0MTU2OCwidXNlck5hbWUiOiLpq5jmiYsiLCJ1c2VyVGFncyI6IjAiLCJ1c2VyVHlwZSI6InVzZXIifQ.EotHOBVvfUuVPYICXizVDsBJAcQ7UJF_fY5CgqVrQcc", nil)
	if err != nil {
		glog.Errorln(err)
	}
	defer c.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("正在等待收消息")
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				glog.Errorln(err)
				continue
			}
			msg := &proto_room.Message{}
			if err := json.Unmarshal(data, msg); err != nil {
				glog.Errorln(err)
			}
			fmt.Println(msg)
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			<-ticker.C
			err := c.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Millisecond*10))
			if err != nil {
				glog.Errorln(err)
				continue
			}
		}
	}()

	wg.Wait()
}
