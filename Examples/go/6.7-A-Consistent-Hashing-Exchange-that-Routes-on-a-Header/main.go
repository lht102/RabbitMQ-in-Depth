package main

import (
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

	if err := ch.ExchangeDeclare(
		"image-storage",
		"x-consistent-hash",
		false,
		false,
		false,
		false,
		amqp.Table{
			"hash-header": "image-hash",
		},
	); err != nil {
		log.Panicf("Failed to declare an exchange: %v", err)
	}
}
