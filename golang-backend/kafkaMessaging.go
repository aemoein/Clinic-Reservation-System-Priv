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
