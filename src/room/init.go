package room

import "github.com/prometheus/client_golang/prometheus"

var (
	Agent Rooms
	Codec codec
	//SenderParallelConnId *senderParallelConnId
)

var (
	// 正在拉取记录的客户端数量
	ClientPullingRecordGauge prometheus.Gauge
	// 正在推送缓存的客户端数量
	ClientPublishCacheGauge prometheus.Gauge

	// 正在队列中的消息
	BroadcastMessageOnQueue prometheus.Gauge
)

func init() {
	Codec.init(10000)

	Agent.Init()
	Sender.Init()
	SyncRecord.Init(100)
	//SenderParallelConnId = &senderParallelConnId{}
	//SenderParallelConnId.Init(100)

	ClientPullingRecordGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "chat_room",
		Name:      "client_pulling_record",
		Help:      "",
	})
	prometheus.MustRegister(ClientPullingRecordGauge)

	ClientPublishCacheGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "chat_room",
		Name:      "client_publish_cache",
		Help:      "",
	})
	prometheus.MustRegister(ClientPublishCacheGauge)

	BroadcastMessageOnQueue = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "chat_room",
		Name:      "broadcast_message_on_queue",
		Help:      "",
	})
	prometheus.MustRegister(BroadcastMessageOnQueue)
}
