package main

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	brokerAddress = "localhost:9092"
	topic         = "clinic_reservation"
)

// Function to produce Kafka message
func produceKafkaMessage(doctorID int, patientID int, operation string) {
	// Create a Kafka producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokerAddress,
	})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		return
	}
	defer p.Close()

	// Create a Kafka message
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(fmt.Sprintf(`{"doctorId": "%d", "patientId": "%d", "Operation": "%s"}`, doctorID, patientID, operation)),
	}

	// Produce the message
	p.ProduceChannel() <- message
}
