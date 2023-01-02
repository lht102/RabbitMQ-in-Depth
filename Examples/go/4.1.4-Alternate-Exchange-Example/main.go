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

	alternateExchangeName := "my-ae"
	if err := ch.ExchangeDeclare(
		alternateExchangeName,
		amqp.ExchangeFanout,
		false,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Panicf("Failed to declare an exchange: %v", err)
	}

	if err := ch.ExchangeDeclare(
		"graphite",
		amqp.ExchangeTopic,
		false,
		false,
		false,
		false,
		amqp.Table{
			"alternate-exchange": alternateExchangeName,
		},
	); err != nil {
		log.Panicf("Failed to declare an exchange: %v", err)
	}

	queueName := "unroutable-messages"
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
		"#",
		alternateExchangeName,
		false,
		nil,
	); err != nil {
		log.Panicf("Failed to bind a queue: %v", err)
	}

	fmt.Print("Queue bound to alternate-exchange")
}
