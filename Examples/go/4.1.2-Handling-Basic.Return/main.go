package main

import (
	"context"
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

	basicReturnCh := make(chan amqp.Return)
	ch.NotifyReturn(basicReturnCh)

	if err := ch.PublishWithContext(
		context.Background(),
		"chapter2-example",
		"server-metrics",
		true,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Timestamp:   time.Now(),
			Type:        "graphite metric",
			Body:        []byte("server.cpu.utilization 25.5 1350884514"),
		},
	); err != nil {
		log.Panicf("Failed to publish a message: %v", err)
	}

	basicReturn := <-basicReturnCh
	fmt.Printf(
		"reply code: %d, reply text: %s, exchange name: %s",
		basicReturn.ReplyCode,
		basicReturn.ReplyText,
		basicReturn.Exchange,
	)
}
