package library

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"data-connector/log"
	"data-connector/utils"

	"github.com/segmentio/kafka-go"
)

const PING_MEM_KEYWORD = "lastTimePingLuffy"
const LUFFY_URL = "https://luffy.`symper`.vn/consumers"
const BOOTSTRAP_BROKER1 = "symperkafla.symper.vn:9092"
const TIMEOUT = 60 //s

type MessageBus struct {
	Topic    string
	Event    string
	Resource interface{}
}

type MessageBusInterface interface {
	PublishBulk() error
	Publish() error
	GetListTopic() map[string]struct{}
	CheckTopicExist() bool
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Xử lý đẩy dữ liệu vào kafka dạng batch ( resource dạng []interface{})
*/
func (messageBus MessageBus) PublishBulk() error {
	// to produce messages
	topic := GetPrefixEnvironment() + messageBus.Topic
	if ok := messageBus.CheckTopicExist(topic); !ok {
		conn, err := kafka.DialLeader(context.Background(), "tcp", BOOTSTRAP_BROKER1, topic, 0)
		if err != nil {
			panic(err.Error())
		}
		defer conn.Close()
	}
	// make a writer that produces to topic-A, using the least-bytes distribution
	w := &kafka.Writer{
		Addr:     kafka.TCP(BOOTSTRAP_BROKER1),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	for i := 0; i < len(messageBus.Resource.([]interface{})); i++ {
		dataPayload := map[string]interface{}{
			"event": messageBus.Event,
			"time":  utils.GetCurrentTimeStamp(),
			"data":  messageBus.Resource.([]interface{})[i],
		}
		dataPayloadJson, _ := json.Marshal(dataPayload)
		err := w.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(messageBus.Event),
			Value: []byte(dataPayloadJson),
		})
		if err != nil {
			fmt.Println(err)
			log.Error("failed to write messages:", map[string]interface{}{
				"trace": log.Trace(),
				"err":   err,
			})
		}
	}
	if err := w.Close(); err != nil {
		fmt.Println(err)
		log.Error("failed to close messages:", map[string]interface{}{
			"trace": log.Trace(),
			"err":   err,
		})
		return err
	}
	return nil

}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Lấy list topic hiện có trên broker
*/
func (messageBus MessageBus) GetListTopic() map[string]struct{} {
	conn, err := kafka.Dial("tcp", BOOTSTRAP_BROKER1)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		panic(err.Error())
	}

	m := map[string]struct{}{}

	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}
	return m
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Kiểm tra xem 1 topic có tồn tại trong list topic hay không, (cần kiểm tra trước khi publish data nếu chưa có topic thì cần tạo topic trước)
*/
func (messageBus MessageBus) CheckTopicExist(topic string) bool {
	l := messageBus.GetListTopic()
	if _, found := l[topic]; found {
		return true
	} else {
		return false
	}
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: publish data vào broker
*/
func (messageBus MessageBus) Publish() error {
	messageBus.Resource = []interface{}{messageBus.Resource}
	err := messageBus.PublishBulk()
	return err
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: hàm subscribe topic từ kafka
*/
func (messageBus MessageBus) SubscribeMultiTopic(topics []string, consumerId, triggerUrl, stopUrl string) error {
	if len(topics) == 0 {
		return errors.New("require topic")
	} else {
		for i := 0; i < len(topics); i++ {
			go GetData(topics[i], consumerId)
		}
	}
	return nil
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Khởi tạo reader lắng nghe dữ liệu từ các topic
*/
func GetData(topic string, consumerId string) {
	t := GetPrefixEnvironment() + topic
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{BOOTSTRAP_BROKER1},
		GroupID:   consumerId,
		Topic:     t,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}
	if err := r.Close(); err != nil {
		fmt.Println("failed to close reader:", err)
	}
}
