package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	brokerAddress = "localhost:9092"
	topic         = "clinic_reservation"
)

func produceKafkaMessage(doctorID int, patientID int, operation string) {
	// Create a struct to represent the message data
	type Message struct {
		DoctorID  int    `json:"doctorId"`
		PatientID int    `json:"patientId"`
		Operation string `json:"operation"`
	}

	// Create a Message instance with the provided data
	message := Message{
		DoctorID:  doctorID,
		PatientID: patientID,
		Operation: operation,
	}

	// Convert the struct to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Connect to Kafka and produce the JSON-formatted message
	conn, err := kafka.DialLeader(context.Background(), "tcp", brokerAddress, topic, 0)
	if err != nil {
		log.Println("Error connecting to Kafka:", err)
		return
	}

	conn.SetWriteDeadline(time.Now().Add(time.Second * 3))
	conn.WriteMessages(kafka.Message{Value: jsonData})
}

func consumeOldKafkaMessages() []string {
	conn, err := kafka.DialLeader(context.Background(), "tcp", brokerAddress, topic, 0)
	if err != nil {
		log.Println("Error connecting to Kafka:", err)
	}

	conn.SetReadDeadline(time.Now().Add(time.Second))

	batch := conn.ReadBatch(1e3, 1e9)
	defer batch.Close()

	var messages []string

	bytes := make([]byte, 1e3)
	for {
		n, err := batch.Read(bytes)
		if err != nil {
			break
		}
		message := string(bytes[:n])
		messages = append(messages, message)
	}

	return messages
}

func consumeKafkaMessages() kafka.Message {
	conn, _ := kafka.DialLeader(context.Background(), "tcp", brokerAddress, topic, 0)
	conn.SetReadDeadline(time.Now().Add(time.Second * 2))

	var lastMessage kafka.Message
	for {
		message, err := conn.ReadMessage(1e3)
		if err != nil {
			//fmt.Println("Error reading Kafka message:", err)
			break
		}

		lastMessage = message
	}

	return lastMessage
}

/*

var (
	brokerAddress = "localhost:9092"
	topic         = "clinic_reservation"
	groupID       = "my-consumer-group"
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
		Value:          []byte(fmt.Sprintf(`{"doctorId": "%d", "patientId": "%d", "operation": "%s"}`, doctorID, patientID, operation)),
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition.Error)
					// Handle the error appropriately
				} else {
					fmt.Printf("Message delivered to topic %s [partition %d] at offset %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			}
		}
	}()

	// Produce the message
	err = p.Produce(message, nil)
	if err != nil {
		fmt.Printf("Failed to produce message: %v\n", err)
		// Handle the error appropriately
	}

	// Wait for any outstanding messages to be delivered and delivery reports to be received
	p.Flush(15 * 1000) // 15-second timeout in milliseconds
}

func consumeKafkaMessages() (*kafka.Consumer, error) {
	// Create a Kafka consumer
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokerAddress,
		"group.id":          groupID,
		"auto.offset.reset": "earliest", // Start consuming from the beginning of the topic if no offset is stored
	})

	if err != nil {
		log.Printf("Failed to create consumer: %v\n", err)
		return nil, err
	}
	defer c.Close()

	// Subscribe to the topic
	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Printf("Error subscribing to topic: %v\n", err)
		return nil, err
	}

	// Consume messages concurrently
	for i := 0; i < 3; i++ { // Adjust the number of goroutines based on your needs
		go func() {
			for {
				msg, err := c.ReadMessage(-1)
				if err == nil {
					fmt.Printf("Received message: %s\n", msg.Value)
					// Process the message as needed
				} else {
					fmt.Printf("Error consuming message: %v\n", err)
				}
			}
		}()
	}

	// Return the Kafka consumer
	return c, err
}
*/
