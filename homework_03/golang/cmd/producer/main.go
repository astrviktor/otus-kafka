package main

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
)

func main() {
	log.Println("Kafka producer start")

	kafkaConfig := kafka.ConfigMap{
		"bootstrap.servers": "localhost:9091",
		"transactional.id":  "transactional-id",
	}

	producer, err := kafka.NewProducer(&kafkaConfig)
	if err != nil {
		log.Println(fmt.Sprintf("NewProducer error: %s", err.Error()))
		return
	}
	defer producer.Close()

	err = producer.InitTransactions(context.Background())
	if err != nil {
		log.Println(fmt.Sprintf("InitTransactions error: %s", err.Error()))
		return
	}

	// Delivery report handler for produced messages
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Delivery message '%s' failed: %v\n", ev.Value, ev.TopicPartition)
				} else {
					log.Printf("Delivered message '%s' to %v\n", ev.Value, ev.TopicPartition)
				}
			}
		}
	}()

	// Transaction 1
	err = producer.BeginTransaction()
	if err != nil {
		log.Println(fmt.Sprintf("BeginTransaction error: %s", err.Error()))
		return
	}

	// Produce messages to topics (asynchronously)
	topic1, topic2 := "topic1", "topic2"
	for _, word := range []string{"message-1", "message-2", "message-3", "message-4", "message-5"} {
		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic1, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)

		if err != nil {
			log.Println(fmt.Sprintf("Produce to topic1 error: %s", err.Error()))
			return
		}

		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic2, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)

		if err != nil {
			log.Println(fmt.Sprintf("Produce to topic2 error: %s", err.Error()))
			return
		}
	}

	// Commit transaction 1
	err = producer.CommitTransaction(context.Background())
	if err != nil {
		log.Println(fmt.Sprintf("CommitTransaction error: %s", err.Error()))
		return
	}

	// Transaction 2
	err = producer.BeginTransaction()
	if err != nil {
		log.Println(fmt.Sprintf("BeginTransaction error: %s", err.Error()))
		return
	}

	// Produce messages to topics (asynchronously)
	for _, word := range []string{"message-6", "message-7"} {
		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic1, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)

		if err != nil {
			log.Println(fmt.Sprintf("Produce to topic1 error: %s", err.Error()))
			return
		}

		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic2, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)

		if err != nil {
			log.Println(fmt.Sprintf("Produce to topic2 error: %s", err.Error()))
			return
		}
	}

	// Abort transaction 2
	err = producer.AbortTransaction(context.Background())
	if err != nil {
		log.Println(fmt.Sprintf("AbortTransaction error: %s", err.Error()))
		return
	}

	// Wait for message deliveries before shutting down
	producer.Flush(15 * 1000)
	log.Println("Kafka producer finish")
}
