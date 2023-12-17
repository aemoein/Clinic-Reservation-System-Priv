package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

var (
	rabbitMQURL   = "amqp://guest:guest@rabbitmqcontainer:5672/"
	exchangeName  = ""
	routingKey    = "clinic_reservation"
	consumerQueue = "clinic_reservation"
)

type Message struct {
	DoctorID  int    `json:"doctorId"`
	PatientID int    `json:"patientId"`
	Operation string `json:"operation"`
}

func produceRabbitMQMessage(doctorID int, patientID int, operation string) error {
	message := Message{
		DoctorID:  doctorID,
		PatientID: patientID,
		Operation: operation,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return fmt.Errorf("error connecting to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("error opening RabbitMQ channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"clinic_reservation",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(q)

	err = ch.Publish(
		exchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	)
	if err != nil {
		return fmt.Errorf("error publishing message to RabbitMQ: %v", err)
	}

	// Log a message indicating that the message was successfully sent
	log.Printf("Message successfully sent to RabbitMQ: DoctorID=%d, PatientID=%d, Operation=%s", doctorID, patientID, operation)

	return nil
}

func consumeRabbitMQMessages() string {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Println("Error connecting to RabbitMQ:", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Println("Error opening RabbitMQ channel:", err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		consumerQueue, // queue
		"",            // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		log.Println("Error consuming messages from RabbitMQ:", err)
	}

	var lastMessage string
	for msg := range msgs {
		lastMessage = string(msg.Body)
	}

	log.Println("consumed messages from RabbitMQ:", lastMessage)
	return lastMessage
}
