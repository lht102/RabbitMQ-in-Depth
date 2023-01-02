package main

import (
	"bytes"
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

	queueName := "test-messages"
	if _, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Panicf("Failed to declare a queue: %v", err)
	}

	for {
		msg, ok, err := ch.Get(queueName, false)
		if err != nil {
			log.Panicf("Failed to get a message: %v", err)
		}

		if !ok {
			continue
		}

		amqputil.PrintMessage(&msg, false)

		if err := msg.Ack(false); err != nil {
			log.Panicf("Failed to acknowledge a message: %v", err)
		}

		if bytes.Equal(msg.Body, []byte("stop")) {
			break
		}
	}
}
