package frontier

import "github.com/prometheus/client_golang/prometheus"

const (
	labelReason = "reason"
)

var (
	// sendDurationsHistogram 统计消息发送耗时分布.
	sendDurationsHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "frontier",
		Subsystem: "message",
		Name:      "send_durations_histogram_seconds",
		Help:      "send message latency distributions.",
	})

	// failedSendCounter 统计发送失败消息.
	failedSendCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "frontier",
		Subsystem: "message",
		Name:      "send_failed_counter",
		Help:      "the count of send message failed.",
	}, []string{labelReason})
)

func init() {
	prometheus.MustRegister(sendDurationsHistogram)
	prometheus.MustRegister(failedSendCounter)
}
