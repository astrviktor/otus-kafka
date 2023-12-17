package main

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	log.Println("Kafka consumer start")

	kafkaConfig := kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(&kafkaConfig)
	if err != nil {
		log.Println(fmt.Sprintf("NewConsumer error: %s", err.Error()))
		return
	}

	err = consumer.SubscribeTopics([]string{"topic1", "topic2"}, nil)
	if err != nil {
		log.Println(fmt.Sprintf("SubscribeTopics error: %s", err.Error()))
		return
	}

	// A signal handler or similar could be used to set this to false to break the loop.
	run := true

	for run {
		msg, err := consumer.ReadMessage(time.Second)
		if err == nil {
			log.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else if !err.(kafka.Error).IsTimeout() {
			// The client will automatically try to recover from all errors.
			// Timeout is not considered an error because it is raised by
			// ReadMessage in absence of messages.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	err = consumer.Close()
	if err != nil {
		log.Println(fmt.Sprintf("SubscribeTopics error: %s", err.Error()))
		return
	}

	log.Println("Kafka consumer finish")
}
