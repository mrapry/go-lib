package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mrapry/go-lib/tracer"

	"github.com/Shopify/sarama"
	"github.com/mrapry/go-lib/golibhelper"
)

// KafkaPublisher kafka
type KafkaPublisher struct {
	producer sarama.SyncProducer
}

// NewKafkaPublisher constructor
func NewKafkaPublisher(brokers []string, cfg *sarama.Config) *KafkaPublisher {

	if len(brokers) == 0 || (len(brokers) == 1 && brokers[0] == "") {
		fmt.Printf(golibhelper.StringYellow("(Kafka publisher: warning, missing kafka broker for publish message. Should be panicked when using kafka publisher.) "))
		return nil
	}

	producer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		fmt.Printf(golibhelper.StringYellow("(Kafka publisher: warning, %v. Should be panicked when using kafka publisher.) "), err)
		return nil
	}

	return &KafkaPublisher{producer}
}

// PublishMessage method
func (p *KafkaPublisher) PublishMessage(ctx context.Context, topic, key string, data interface{}) (err error) {
	opName := "kafka:publish_message"

	var payload []byte

	switch d := data.(type) {
	case string:
		payload = []byte(d)
	case []byte:
		payload = d
	default:
		payload, _ = json.Marshal(data)
	}

	tracer.WithTraceFunc(ctx, opName, func(c context.Context, tag map[string]interface{}) {
		defer func() {
			// set tracer tag
			tag["topic"] = topic
			tag["key"] = key
			tag["message"] = string(payload)

			msg := &sarama.ProducerMessage{
				Topic:     topic,
				Key:       sarama.ByteEncoder([]byte(key)),
				Value:     sarama.ByteEncoder(payload),
				Timestamp: time.Now(),
			}
			_, _, err = p.producer.SendMessage(msg)
		}()
	})

	return
}
