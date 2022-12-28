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

	if err := ch.Qos(
		10,
		0,
		false,
	); err != nil {
		log.Panicf("Failed to set QoS: %v", err)
	}

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

	unacknowledged := 0

	for delivery := range deliveryCh {
		amqputil.PrintMessage(&delivery, false)
		unacknowledged++

		if unacknowledged == 10 {
			if err := delivery.Ack(true); err != nil {
				log.Panicf("Failed to acknowledge a message: %v", err)
			}

			unacknowledged = 0
		}
	}
}
