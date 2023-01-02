package main

import (
	"fmt"
	"log"

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

	for i := 0; i < 4; i++ {
		queueName := fmt.Sprintf("server%d", i)
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

		if err := ch.QueueBind(
			queueName,
			"10",
			"image-storage",
			false,
			nil,
		); err != nil {
			log.Panicf("Failed to bind a queue: %v", err)
		}
	}
}
