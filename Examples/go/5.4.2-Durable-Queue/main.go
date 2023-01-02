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

	if _, err := ch.QueueDeclare(
		"durable-queue",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Panicf("Failed to declare a queue: %v", err)
	}

	fmt.Print("Queue declared")
}
