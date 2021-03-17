package main

import (
	"log"
	"os"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func main() {
	kafkaServer := readFromENV("KAFKA_BROKER", "localhost:9092")
	kafkaTopic := readFromENV("KAFKA_TOPIC", "testTopic")
	kafkaGroup := readFromENV("GROUP_ID", "myGroup")

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaServer,
		"group.id":          kafkaGroup,
		"auto.offset.reset": "earliest",

	})

	if err != nil {
		log.Fatalln(err)
	}

	c.SubscribeTopics([]string{kafkaTopic}, nil)
	defer c.Close()

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			log.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else {
			// The client will automatically try to recover from all errors.
			log.Fatalf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}

func readFromENV(key, defaultVal string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	return value
}
