package main

import (
	"encoding/json"
	"github.com/selcux/data-pipeline-sample/ViewProducerApp/internal"
	"github.com/selcux/data-pipeline-sample/ViewProducerApp/model"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	source := os.Getenv("FILE_SOURCE")
	fileOps := internal.NewFileOps(source)

	lineCh := make(chan string)
	errorCh := make(chan error)
	closeCh := make(chan struct{})
	go fileOps.ReadInterval(1*time.Second, lineCh, errorCh, closeCh)

	kafkaServer := readFromENV("KAFKA_BROKER", "localhost:9092")
	kafkaTopic := readFromENV("KAFKA_TOPIC", "product-view")

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaServer})
	if err != nil {
		panic(err)
	}
	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					log.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the app
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	for {
		select {
		case line := <-lineCh:
			data := Serialize(line)
			err := p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &kafkaTopic, Partition: kafka.PartitionAny},
				Value:          data,
			}, nil)
			if err != nil {
				log.Fatalln("p.Produce", err)
			}
		case err := <-errorCh:
			log.Fatalln(err)
		case <-closeCh:
			log.Println("Closing")
			// Wait for message deliveries before shutting down
			p.Flush(15 * 1000)
			return
		case <-quit:
			log.Println("Quitting")
			return
		}
	}
}

func Serialize(data string) []byte {
	productView := model.ProductView{}
	err := json.Unmarshal([]byte(data), &productView)
	if err != nil {
		log.Fatalln("Unmarshal", err, data)
	}
	productView.Timestamp = time.Now()

	byteData, err := json.Marshal(productView)
	if err != nil {
		log.Fatalln("Marshal", err, data)
	}

	return byteData
}

func readFromENV(key, defaultVal string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	return value
}
