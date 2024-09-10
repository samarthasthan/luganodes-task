package main

import (
	"github.com/samarthasthan/luganodes-task/internal/store/database"
	consumer "github.com/samarthasthan/luganodes-task/internal/store/kafka"
	"github.com/samarthasthan/luganodes-task/pkg/env"
)

var (
	KAFKA_PORT          string
	KAFKA_HOST          string
	MYSQL_PORT          string
	MYSQL_ROOT_PASSWORD string
	MYSQL_HOST          string
	REDIS_PORT          string
	REDIS_HOST          string
)

func init() {
	KAFKA_PORT = env.GetEnv("KAFKA_PORT", "9092")
	KAFKA_HOST = env.GetEnv("KAFKA_HOST", "localhost")
	MYSQL_PORT = env.GetEnv("MYSQL_PORT", "3306")
	MYSQL_ROOT_PASSWORD = env.GetEnv("MYSQL_ROOT_PASSWORD", "password")
	MYSQL_HOST = env.GetEnv("MYSQL_HOST", "localhost")
}

func main() {
	// Create mysql database
	sql := database.NewMySQL()
	sql.Connect("root:" + MYSQL_ROOT_PASSWORD + "@tcp(" + MYSQL_HOST + ":" + MYSQL_PORT + ")/luganodes")
	defer sql.Close()

	// Create kafka consumer
	kafka := consumer.NewKafkaConsumer(KAFKA_HOST+":"+KAFKA_PORT, sql)

	// Consume messages from kafka and store in mysql and redis
	kafka.Consume()
}
