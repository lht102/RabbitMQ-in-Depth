package main

import (
	"context"
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

	exchangeName := "chapter4-example"
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

	if err := ch.Confirm(false); err != nil {
		log.Panicf("Failed to enable publisher confirms: %v", err)
	}

	confirmation, err := ch.PublishWithDeferredConfirmWithContext(
		context.Background(),
		exchangeName,
		"important.message",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Type:        "very important",
			Body:        []byte("This is an important message"),
		},
	)
	if err != nil {
		log.Panicf("Failed to publish a message: %v", err)
	}

	if confirmation.Wait() {
		fmt.Print("The message was confirmed")
	}
}
