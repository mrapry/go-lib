package kafkaworker

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/mrapry/go-lib/codebase/factory"
	"github.com/mrapry/go-lib/codebase/factory/types"
	"github.com/mrapry/go-lib/config"
	"github.com/mrapry/go-lib/golibhelper"
	"github.com/mrapry/go-lib/logger"
	"github.com/mrapry/go-lib/tracer"
)

type kafkaWorker struct {
	engine  sarama.ConsumerGroup
	service factory.ServiceFactory
}

// NewWorker create new kafka consumer
func NewWorker(service factory.ServiceFactory) factory.AppServerFactory {
	// init kafka consumer
	kafkaConsumer, err := sarama.NewConsumerGroup(
		config.BaseEnv().Kafka.Brokers,
		config.BaseEnv().Kafka.ConsumerGroup,
		service.GetDependency().GetBroker().GetConfig(),
	)
	if err != nil {
		log.Panicf("Error creating kafka consumer group client: %v", err)
	}

	return &kafkaWorker{
		engine:  kafkaConsumer,
		service: service,
	}
}

func (h *kafkaWorker) Serve() {

	topicInfo := make(map[string]string)
	handlerFuncs := make(map[string]types.WorkerHandlerFunc)
	for _, m := range h.service.GetModules() {
		if h := m.WorkerHandler(types.Kafka); h != nil {
			for topic, handlerFunc := range h.MountHandlers() {
				handlerFuncs[topic] = handlerFunc
				topicInfo[topic] = string(m.Name())
			}
		}
	}

	consumer := kafkaConsumer{
		handlerFuncs: handlerFuncs,
	}

	var consumeTopics []string
	for topic, moduleName := range topicInfo {
		fmt.Println(golibhelper.StringYellow(fmt.Sprintf("[KAFKA-CONSUMER] (topic): %-8s  (consumed by module)--> [%s]", topic, moduleName)))
		consumeTopics = append(consumeTopics, topic)
	}
	fmt.Printf("\x1b[34;1m⇨ Kafka consumer is active. Brokers: " + strings.Join(config.BaseEnv().Kafka.Brokers, ", ") + "\x1b[0m\n\n")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := h.engine.Consume(ctx, consumeTopics, &consumer); err != nil {
		log.Panicf("Error from kafka consumer: %v", err)
	}
}

func (h *kafkaWorker) Shutdown(ctx context.Context) {
	deferFunc := logger.LogWithDefer("Stopping Kafka consumer worker...")
	defer deferFunc()

	h.engine.Close()
}

// kafkaConsumer represents a Sarama consumer group consumer
type kafkaConsumer struct {
	handlerFuncs map[string]types.WorkerHandlerFunc // mapping topic to handler func in delivery layer
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *kafkaConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *kafkaConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *kafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)

		tracer.WithTraceFunc(session.Context(), "KafkaConsumer", func(ctx context.Context, tags map[string]interface{}) {
			defer func() {
				if r := recover(); r != nil {
					tracer.SetError(ctx, fmt.Errorf("%v", r))
				} else {
					session.MarkMessage(message, "")
				}
			}()

			tags["topic"] = message.Topic
			tags["key"] = string(message.Key)
			tags["value"] = string(message.Value)

			handlerFunc := c.handlerFuncs[message.Topic]
			if err := handlerFunc(ctx, message.Value); err != nil {
				panic(err)
			}
		})
	}

	return nil
}
