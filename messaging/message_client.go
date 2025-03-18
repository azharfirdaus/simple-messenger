package messaging

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"
)

type ChatQueueClient interface {
	Produce(userID, data *string) error
	Consume()
}

type KafkaChatQueueClientImpl struct {
	broker *string
	topic  *string
}

func NewKafkaMessageQueueClientImpl(broker string, topic string) *KafkaChatQueueClientImpl {
	return &KafkaChatQueueClientImpl{
		broker: &broker,
		topic:  &topic,
	}
}

func (k KafkaChatQueueClientImpl) Produce(ID *uint64, data *string) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = 0

	producer, err := sarama.NewSyncProducer([]string{*k.broker}, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
		return err
	}

	defer producer.Close()

	protobufChat := &ProtobufChat{ID: *ID, Data: *data}
	protobufChatBytes, err := proto.Marshal(protobufChat)
	if err != nil {
		return fmt.Errorf("failed to marshal: %v", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: *k.topic,
		Value: sarama.ByteEncoder(protobufChatBytes),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatalf("Failed to send message to Kafka: %v", err)
		return err
	}

	// Log the success
	log.Printf("Message sent successfully to topic '%s' (partition %d, offset %d)\n", *k.topic, partition, offset)
	return nil
}

func (k KafkaChatQueueClientImpl) Consume(consumerGroupHandler *ConsumerGroupHandler) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup([]string{*k.broker}, "my-group", config)
	if err != nil {
		log.Fatalf("Failed to create consumer group: %v", err)
	}
	defer consumerGroup.Close()

	ctx := context.Background()
	log.Printf("Starting consumer for topic: %s\n", *k.topic)
	for {
		err := consumerGroup.Consume(ctx, []string{*k.topic}, consumerGroupHandler)
		if err != nil {
			log.Fatalf("Error from consumer: %v", err)
		}
	}
}

type ConsumerGroupHandler struct {
	Handler func(topic *string, partition *int32, offset *int64, key, value *[]byte)
}

func (*ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (*ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		h.Handler(&message.Topic, &message.Partition, &message.Offset, &message.Key, &message.Value)
		session.MarkMessage(message, "") // Mark message as processed
	}
	return nil
}
