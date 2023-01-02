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

	if err := ch.Tx(); err != nil {
		log.Panicf("Failed to select standard transaction mode: %v", err)
	}

	if err := ch.PublishWithContext(
		context.Background(),
		"chapter4-example",
		"important.message",
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp.Persistent,
			Type:         "important",
			Body:         []byte("This is an important message"),
		},
	); err != nil {
		log.Panicf("Failed to publish a message: %v", err)
	}

	if err := ch.TxCommit(); err != nil {
		log.Panicf("Failed to commit the current transaction: %v", err)
	}

	fmt.Print("Transaction committed")
}
