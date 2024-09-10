package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/labstack/gommon/log"
	"github.com/samarthasthan/luganodes-task/internal/fetcher/models"
	"github.com/samarthasthan/luganodes-task/internal/store/database"
	"github.com/samarthasthan/luganodes-task/internal/store/database/mysql/sqlc"
	"github.com/samarthasthan/luganodes-task/pkg/logger"
)

type Consumer interface{}

type KafkaConsumer struct {
	Consumer *kafka.Consumer
	mysql    database.Database
	log      *logger.Logger
}

func NewKafkaConsumer(kafkaURL string, db database.Database, logger *logger.Logger) *KafkaConsumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaURL,
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		fmt.Println(err)
	}
	return &KafkaConsumer{
		Consumer: c,
		mysql:    db,
		log:      logger,
	}
}

func (c *KafkaConsumer) Consume() {
	c.Consumer.SubscribeTopics([]string{"fetcher"}, nil)
	for {
		msg, err := c.Consumer.ReadMessage(time.Second)
		if err == nil {

			// Consume Tweet fromm Kafka
			var deposit models.Deposit
			json.Unmarshal(msg.Value, &deposit) // Marshal Kafka Message to Tweet Struct

			log.Infof("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))

			// Type assertion for MySQL
			mysql, ok := c.mysql.(*database.MySQL)
			if !ok {
				log.Errorf("Error converting to MySQL: %v", err)
			}

			ctx := context.Background()

			tx, err := mysql.DB.BeginTx(ctx, nil)
			if err != nil {
				log.Errorf("Error starting transaction: %v", err)
			}

			err = mysql.Queries.InsertDeposit(ctx, sqlc.InsertDepositParams{
				Blocknumber:    int32(deposit.BlockNumber),
				Blocktimestamp: int32(deposit.BlockTimestamp),
				Fee:            deposit.Fee,
				Hash:           deposit.Hash,
				Pubkey:         deposit.Pubkey,
			})

			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}

		} else if !err.(kafka.Error).IsTimeout() {
			log.Errorf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
