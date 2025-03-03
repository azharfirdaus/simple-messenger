package messaging

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

type MessageQueueClient interface {
	Produce(message string) error
	Consume()
}

type KafkaMessageQueueClientImpl struct {
	broker *string
	topic  *string
}

func NewKafkaMessageQueueClientImpl(broker string, topic string) *KafkaMessageQueueClientImpl {
	return &KafkaMessageQueueClientImpl{
		broker: &broker,
		topic:  &topic,
	}
}

func (k *KafkaMessageQueueClientImpl) Produce(message string) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = 0

	producer, err := sarama.NewSyncProducer([]string{*k.broker}, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
		return err
	}

	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: *k.topic,
		Value: sarama.StringEncoder(message),
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

func (k *KafkaMessageQueueClientImpl) Consume() {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup([]string{*k.broker}, "my-group", config)
	if err != nil {
		log.Fatalf("Failed to create consumer group: %v", err)
	}
	defer consumerGroup.Close()

	handler := &ConsumerGroupHandler{}

	ctx := context.Background()

	// Start consuming messages
	log.Printf("Starting consumer for topic: %s\n", *k.topic)
	for {
		// Consume messages
		err := consumerGroup.Consume(ctx, []string{*k.topic}, handler)
		if err != nil {
			log.Fatalf("Error from consumer: %v", err)
		}
	}
}

type ConsumerGroupHandler struct{}

func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Received message: %s (topic: %s, partition: %d, offset: %d)\n",
			string(message.Value), message.Topic, message.Partition, message.Offset)
		session.MarkMessage(message, "") // Mark message as processed
	}
	return nil
}
