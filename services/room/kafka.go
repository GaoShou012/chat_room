package room

import (
	"fmt"
	"github.com/Shopify/sarama"
	"strconv"
	"time"
	proto_room "wchatv1/proto/room"
)

type BroadcastToFrontierByKafka struct {
	topic    string
	producer sarama.AsyncProducer
	bucket   chan *proto_room.Message
	//counter  int
}

func (b *BroadcastToFrontierByKafka) Init(kafkaAddr []string, topic string) error {
	b.topic = topic
	b.bucket = make(chan *proto_room.Message, 100000)

	conf := sarama.NewConfig()
	conf.Producer.Retry.Max = 5
	conf.Producer.RequiredAcks = sarama.WaitForAll
	producer, err := sarama.NewAsyncProducer(kafkaAddr, conf)
	if err != nil {
		return err
	}
	b.producer = producer

	//go func() {
	//	ticker := time.NewTicker(time.Second)
	//	for {
	//		now := <-ticker.C
	//		fmt.Println(b.counter, now.Unix())
	//	}
	//}()

	go func() {
		defer producer.Close()
		for {
			message := <-b.bucket
			j, err := Codec.EncodeMessage(message)
			if err != nil {
				fmt.Println(err)
				continue
			}
			strTime := strconv.Itoa(int(time.Now().Unix()))
			m := &sarama.ProducerMessage{
				Topic:     b.topic,
				Key:       sarama.StringEncoder(strTime),
				Value:     sarama.StringEncoder(j),
				Headers:   nil,
				Metadata:  nil,
				Offset:    0,
				Partition: 0,
				Timestamp: time.Time{},
			}
			b.producer.Input() <- m
			//b.counter++
		}
	}()

	return nil
}

func (b *BroadcastToFrontierByKafka) Broadcast() chan<- *proto_room.Message {
	return b.bucket
}
func (b *BroadcastToFrontierByKafka) Bucket() chan *proto_room.Message {
	return b.bucket
}
