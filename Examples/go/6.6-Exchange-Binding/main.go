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

	topicExchangeName := "events"
	if err := ch.ExchangeDeclare(
		topicExchangeName,
		amqp.ExchangeTopic,
		false,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Panicf("Failed to declare an exchange: %v", err)
	}

	consistentHashExchangeName := "distributed-events"
	if err := ch.ExchangeDeclare(
		consistentHashExchangeName,
		"x-consistent-hash",
		false,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Panicf("Failed to declare an exchange: %v", err)
	}

	if err := ch.ExchangeBind(
		consistentHashExchangeName,
		"#",
		topicExchangeName,
		false,
		nil,
	); err != nil {
		log.Panicf("Failed to bind an exchange: %v", err)
	}
}
