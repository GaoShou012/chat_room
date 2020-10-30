package room

import (
	"container/list"
	"context"
	"fmt"
	"github.com/GaoShou012/frontier"
	"github.com/golang/glog"
	"os"
	"runtime"
	"strings"
	"sync"
	"wchatv1/config"
	proto_room "wchatv1/proto/room"
	"wchatv1/utils"
)

type onMessageJob struct {
	conn    frontier.Conn
	message []byte
}

type Rooms struct {
	mutex   sync.Mutex
	rooms   map[Code]*Room
	clients map[frontier.Conn]*Client

	BroadcastCache chan []byte
	onMessageJobs  chan *onMessageJob
}

// 房间代理初始化
func (rs *Rooms) Init() {
	// 房间
	rs.rooms = make(map[Code]*Room)
	// 消息广播
	rs.BroadcastCache = make(chan []byte, 100000)
	// 连接的客户端
	rs.clients = make(map[frontier.Conn]*Client)

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				message := <-rs.BroadcastCache
				if msg, err := Codec.Decode(message); err == nil {
					rs.Broadcast(msg)
				} else {
					//glog.Errorln(err)
				}
			}
		}()
	}

	rs.onMessageJobs = make(chan *onMessageJob, 100000)
	for i := 0; i < runtime.NumCPU()*30; i++ {
		go func() {
			for {
				job := <-rs.onMessageJobs
				conn, message := job.conn, job.message
				ctx := conn.GetContext().(*ConnContext)
				msg, err := Codec.Decode(message)
				if err != nil {
					fmt.Println(message)
					glog.Errorln(err)
					continue
				}
				msg.From = &proto_room.From{
					Path: fmt.Sprintf("%s.%d.%s.%s.%d", conn.Frontier().Id, conn.GetId(), ctx.TenantCode, ctx.RoomCode, ctx.UserId),
					PassToken: &proto_room.PassToken{
						TenantCode: ctx.TenantCode,
						RoomCode:   ctx.RoomCode,
						UserType:   ctx.UserType,
						UserId:     ctx.UserId,
						UserName:   ctx.UserName,
						UserThumb:  ctx.UserThumb,
						UserTags:   ctx.UserTags,
					},
				}

				// 请求房间服务
				broadcastReq := &proto_room.BroadcastReq{Message: msg}
				serviceClient := config.RoomServiceConfig.ServiceClient()
				rsp, err := serviceClient.Broadcast(context.TODO(), broadcastReq, utils.Micro.DefaultCliCallOptions)
				if err != nil {
					glog.Errorln(err)
					continue
				}

				bts, err := Codec.Encode(rsp.Message)
				if err != nil {
					glog.Errorln(err)
					continue
				}
				conn.Sender(bts)
			}
		}()
	}
}

// 获取当前客户端数量
func (rs *Rooms) GetClients() int {
	return len(rs.clients)
}

// 缓存线 切换到 直联
func (rs *Rooms) ClientToDirectMode(cli *Client) {
	ctx := cli.conn.GetContext().(*ConnContext)
	roomCode := RoomCode(ctx)

	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	r, ok := rs.rooms[roomCode]
	if !ok {
		glog.Errorln("The room is not exists", roomCode)
		return
	}
	r.ClientToDirectMode(cli)
}

// 连接加入房间
func (rs *Rooms) Join(conn frontier.Conn, passToken *proto_room.PassToken, lastMessageId string) {
	// 绑定上下文
	ctx := &ConnContext{}
	conn.SetContext(ctx)
	ctx.TenantCode = passToken.TenantCode
	ctx.RoomCode = passToken.RoomCode
	ctx.UserType = passToken.UserType
	ctx.UserId = passToken.UserId
	ctx.UserName = passToken.UserName
	ctx.UserThumb = passToken.UserThumb
	ctx.UserTags = passToken.UserTags

	// 生成client
	cli := &Client{
		state:               0,
		conn:                conn,
		isClose:             false,
		isCache:             false,
		isSyncDone:          false,
		cache:               list.New(),
		record:              list.New(),
		lastMessageId:       lastMessageId,
		isPullingRecordDone: false,
		isPublishRecordDone: false,
		isPublishCacheDone:  false,
	}
	roomCode := RoomCode(ctx)

	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	// create room
	r, ok := rs.rooms[roomCode]
	if !ok {
		r = &Room{}
		r.Init(ctx.TenantCode, ctx.RoomCode, 100000, 100000)
		rs.rooms[roomCode] = r
	}

	// save client
	rs.clients[conn] = cli

	// join to the room
	r.Join(cli)
}

// 连接离开房间
func (rs *Rooms) Leave(conn frontier.Conn) {
	ctx := conn.GetContext().(*ConnContext)
	roomCode := RoomCode(ctx)

	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	cli, ok := rs.clients[conn]
	if ok {
		delete(rs.clients, conn)
	} else {
		return
	}

	// 离开房间
	// 如果房间没有客户端，关闭房间
	r, ok := rs.rooms[roomCode]
	if !ok {
		return
	}
	r.Leave(cli)
	if r.ClientsCounter() == 0 {
		delete(rs.rooms, roomCode)
		go r.Close()
	}
}

// 服务器收到新的消息
// 路由到房间服务，检查消息
func (rs *Rooms) OnMessage(conn frontier.Conn, message []byte) {
	job := &onMessageJob{
		conn:    conn,
		message: message,
	}
	rs.onMessageJobs <- job
	//ctx := conn.GetContext().(*ConnContext)
	//msg, err := Codec.Decode(message)
	//if err != nil {
	//	return err
	//}
	//msg.From = &proto_room.From{
	//	Path: fmt.Sprintf("%d.%d.%s.%s.%d", conn.GetFid(), conn.GetId(), ctx.TenantCode, ctx.RoomCode, ctx.UserId),
	//	PassToken: &proto_room.PassToken{
	//		TenantCode: ctx.TenantCode,
	//		RoomCode:   ctx.RoomCode,
	//		UserType:   ctx.UserType,
	//		UserId:     ctx.UserId,
	//		UserName:   ctx.UserName,
	//		UserThumb:  ctx.UserThumb,
	//		UserTags:   ctx.UserTags,
	//	},
	//}
	//
	//// 请求房间服务
	//broadcastReq := &proto_room.BroadcastReq{Message: msg}
	//var opss client.CallOption = func(o *client.CallOptions) {
	//	o.RequestTimeout = time.Second * 10
	//	o.DialTimeout = time.Second * 10
	//}
	//serviceClient := config.RoomServiceConfig.ServiceClient()
	//rsp, err := serviceClient.Broadcast(context.TODO(), broadcastReq, opss)
	//if err != nil {
	//	return err
	//}
	//
	//bts, err := Codec.Encode(rsp.Message)
	//if err != nil {
	//	return err
	//}
	//return conn.Writer(bts)
}

// 广播消息
// 把消息，广播到房间里的客户端
func (rs *Rooms) Broadcast(message *proto_room.Message) {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	if message.To == nil {
		glog.Errorln("message.To is nil", message)
		return
	}

	// 拆分路径，组合出房间代码
	arr := strings.Split(message.To.Path, ".")
	if len(arr) != 2 {
		glog.Errorln("message.To.Path 格式错误", message)
		return
	}
	roomCode := fmt.Sprintf("%s.%s", arr[0], arr[1])

	// 查询房间是否存在
	r, ok := rs.rooms[Code(roomCode)]
	if !ok {
		return
	}

	if message == nil {
		fmt.Println("rooms message is nil")
		os.Exit(1)
	}
	r.Broadcast(message)
}
