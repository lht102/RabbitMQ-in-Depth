package main

import (
	"crypto/md5" //nolint: gosec
	"fmt"
	"log"
	"os"

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

	queueName := fmt.Sprintf("hashing-worker-%d", os.Getpid())
	if _, err := ch.QueueDeclare(
		queueName,
		false,
		true,
		true,
		false,
		nil,
	); err != nil {
		log.Panicf("Failed to declare a queue: %v", err)
	}

	fmt.Println("Worker queue declared")

	if err := ch.QueueBind(
		queueName,
		"",
		"fanout-rpc-requests",
		false,
		nil,
	); err != nil {
		log.Panicf("Failed to bind a queue: %v", err)
	}

	fmt.Println("Worker queue bound")

	deliveryCh, err := ch.Consume(
		queueName,
		"",
		false,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Panicf("Failed to register a consumer: %v", err)
	}

	for msg := range deliveryCh {
		hashObj := md5.Sum(msg.Body) //nolint: gosec

		fmt.Printf(
			"Image with correlation-id of %s has a hash of %x\n",
			msg.CorrelationId,
			hashObj,
		)

		if err := msg.Ack(false); err != nil {
			log.Panicf("Failed to acknowledge a message: %v", err)
		}
	}
}
