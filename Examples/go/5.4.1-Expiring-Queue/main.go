package main

import (
	"fmt"
	"log"
	"time"

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

	queueName := "expiring-queue"
	args := amqp.Table{
		"x-expires": 1000,
	}

	if _, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		args,
	); err != nil {
		log.Panicf("Failed to declare a queue: %v", err)
	}

	if _, err := ch.QueueDeclarePassive(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Panicf("Failed to declare a queue: %v", err)
	}

	time.Sleep(2 * time.Second)

	if _, err := ch.QueueDeclarePassive(
		queueName,
		false,
		false,
		false,
		false,
		args,
	); err != nil {
		fmt.Printf("The queue no longer exists: %v", err)
	}
}
