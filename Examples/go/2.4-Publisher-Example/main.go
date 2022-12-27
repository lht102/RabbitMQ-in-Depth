package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
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

	exchangeName := "chapter2-example"
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

	queueName := "example"
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

	routingKey := "example-routing-key"
	if err := ch.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false,
		nil,
	); err != nil {
		log.Panicf("Failed to bind a queue: %v", err)
	}

	for i := 0; i < 10; i++ {
		if err := ch.PublishWithContext(
			context.Background(),
			exchangeName,
			routingKey,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				MessageId:   uuid.NewString(),
				Timestamp:   time.Now(),
				Body:        []byte(fmt.Sprintf("Text message#%d", i)),
			},
		); err != nil {
			log.Panicf("Failed to publish a message: %v", err)
		}
	}
}
