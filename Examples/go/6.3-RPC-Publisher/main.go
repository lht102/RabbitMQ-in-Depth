package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
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

	queueName := fmt.Sprintf("response-queue-%d", os.Getpid())
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

	fmt.Println("Response queue declared")

	if err := ch.QueueBind(
		queueName,
		queueName,
		"rpc-replies",
		false,
		nil,
	); err != nil {
		log.Panicf("Failed to bind a queue: %v", err)
	}

	fmt.Println("Response queue bound")

	images, err := ch6.GetImages()
	if err != nil {
		log.Panicf("Failed to get images: %v", err)
	}

	for imgID, filename := range images {
		fmt.Printf("Sending request for image #%d: %s\n", imgID, filename)

		body, err := os.ReadFile(filename)
		if err != nil {
			log.Panicf("Failed to read image: %v", err)
		}

		mimeType, err := ch6.MimeType(filename)
		if err != nil {
			log.Panicf("Failed to read mime type: %v", err)
		}

		if err := ch.PublishWithContext(
			context.Background(),
			"topic-rpc-requests",
			"image.new.profile",
			false,
			false,
			amqp.Publishing{
				ContentType:   mimeType,
				CorrelationId: strconv.Itoa(imgID),
				ReplyTo:       queueName,
				MessageId:     uuid.NewString(),
				Timestamp:     time.Now(),
				Body:          body,
			},
		); err != nil {
			log.Panicf("Failed to publish a message: %v", err)
		}

		var message *amqp.Delivery
		for message == nil {
			delivery, ok, err := ch.Get(queueName, false)
			if err != nil {
				log.Panicf("Failed to get a message: %v", err)
			}

			if ok {
				message = &delivery
			}
		}

		if err := message.Ack(false); err != nil {
			log.Panicf("Failed to acknowledge a message: %v", err)
		}

		firstPublishTime, ok := message.Headers["first_publish"].(time.Time)
		if ok {
			duration := time.Since(firstPublishTime)
			fmt.Printf(
				"Facial detection RPC call for image %s total duration: %s\n",
				message.CorrelationId,
				duration,
			)
		}

		if err := ch6.DisplayImage(message.Body, message.ContentType); err != nil {
			log.Panicf("Failed to display image: %v", err)
		}
	}

	fmt.Printf("RPC requests processed")
}
