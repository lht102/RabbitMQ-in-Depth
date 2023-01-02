package main

import (
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

		if err := delivery.Ack(false); err != nil {
			log.Panicf("Failed to acknowledge a message: %v", err)
		}
	}
}
