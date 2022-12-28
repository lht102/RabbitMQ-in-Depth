package main

import (
	"fmt"
	"log"

	"github.com/lht102/RabbitMQ-in-Depth/Examples/go/amqputil"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Panicf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	deliveryCh, err := ch.Consume(
		"test-messages",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Panicf("Failed to register a consumer: %v", err)
	}

	for delivery := range deliveryCh {
		amqputil.PrintMessage(&delivery, false)

		fmt.Printf("Redelivered: %t\n", delivery.Redelivered)

		if err := delivery.Reject(true); err != nil {
			log.Panicf("Failed to reject a message: %v", err)
		}
	}
}
