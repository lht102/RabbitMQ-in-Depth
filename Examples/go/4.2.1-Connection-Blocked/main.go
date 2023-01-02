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

	connectionBlockedCh := make(chan amqp.Blocking)
	conn.NotifyBlocked(connectionBlockedCh)

	blocked := <-connectionBlockedCh
	if blocked.Active {
		fmt.Printf("Connection blocked: %s", blocked.Reason)
	} else {
		fmt.Printf("Connection unblocked")
	}
}
