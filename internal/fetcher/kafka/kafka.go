package kafka

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/samarthasthan/luganodes-task/internal/fetcher/models"
)

type DataProducer interface {
	Produce(deposit *models.Deposit)
}

type KafkaProducer struct {
	Producer *kafka.Producer
}

func NewKafkaProducer(kafkaURL string) *KafkaProducer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaURL})
	if err != nil {
		panic(err)
	}
	return &KafkaProducer{
		Producer: p,
	}
}

func (k *KafkaProducer) Produce(deposit *models.Deposit) {
	topic := "fetcher"
	data, err := json.Marshal(deposit)
	if err != nil {
		log.Fatalln(err)
	}
	k.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(data),
	}, nil)
}
