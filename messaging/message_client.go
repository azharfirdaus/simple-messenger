package messaging

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type MessageQueueClient interface {
	Produce(message string) error
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
	fmt.Printf("Message sent successfully to topic '%s' (partition %d, offset %d)\n", *k.topic, partition, offset)
	return nil
}
