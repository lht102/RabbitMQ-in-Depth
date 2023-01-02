package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/lht102/RabbitMQ-in-Depth/Examples/go/ch6"
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

	queueName := fmt.Sprintf("rpc-worker-%d", os.Getpid())
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
		"detect-faces",
		"direct-rpc-requests",
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
		duration := time.Since(msg.Timestamp)
		fmt.Printf(
			"Received RPC request published %.2f seconds ago\n",
			float64(duration)/float64(time.Second),
		)

		tempFile, err := ch6.WriteTempFile(
			msg.Body,
			msg.ContentType,
		)
		if err != nil {
			log.Panicf("Failed to write temp file: %v", err)
		}

		resultFile, err := ch6.Faces(tempFile)
		if err != nil {
			log.Panicf("Failed to detect faces: %v", err)
		}

		body, err := os.ReadFile(resultFile)
		if err != nil {
			log.Panicf("Failed to read result file: %v", err)
		}

		if err := os.Remove(tempFile); err != nil {
			log.Panicf("Failed to remove temp file: %v", err)
		}

		if err := os.Remove(resultFile); err != nil {
			log.Panicf("Failed to remove result file: %v", err)
		}

		if err := ch.PublishWithContext(
			context.Background(),
			"rpc-replies",
			msg.ReplyTo,
			false,
			false,
			amqp.Publishing{
				Headers: amqp.Table{
					"first_publish": msg.Timestamp,
				},
				ContentType:   msg.ContentType,
				CorrelationId: msg.CorrelationId,
				MessageId:     uuid.NewString(),
				Timestamp:     time.Now(),
				AppId:         "Chapter 6 Listing 2 Consumer",
				Body:          body,
			},
		); err != nil {
			log.Panicf("Failed to publish a message: %v", err)
		}

		if err := msg.Ack(false); err != nil {
			log.Panicf("Failed to acknowledge a message: %v", err)
		}
	}
}
