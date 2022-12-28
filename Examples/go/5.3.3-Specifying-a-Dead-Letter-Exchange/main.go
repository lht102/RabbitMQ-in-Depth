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

	exchangeName := "rejected-messages"
	if err := ch.ExchangeDeclare(
		exchangeName,
		amqp.ExchangeDirect,
		false,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Panicf("Failed to declare an exchange: %v", err)
	}

	if _, err := ch.QueueDeclare(
		"dlx-example",
		false,
		false,
		false,
		false,
		amqp.Table{
			"dead_letter_exchange": exchangeName,
		},
	); err != nil {
		log.Panicf("Failed to declare a queue: %v", err)
	}
}
