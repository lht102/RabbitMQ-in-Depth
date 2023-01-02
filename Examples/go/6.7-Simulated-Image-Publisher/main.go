package main

import (
	"context"
	"crypto/md5" //nolint: gosec
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

	for i := 0; i < 100000; i++ {
		hashValue := md5.Sum([]byte(fmt.Sprintf("%s:%d", time.Now(), i))) //nolint: gosec

		if err := ch.PublishWithContext(
			context.Background(),
			"image-storage",
			"",
			false,
			false,
			amqp.Publishing{
				Headers: amqp.Table{
					"image-hash": fmt.Sprintf("%x", hashValue),
				},
				MessageId: uuid.NewString(),
				Timestamp: time.Now(),
				Body:      []byte(fmt.Sprintf("Image # %d", i)),
			},
		); err != nil {
			log.Panicf("Failed to publish a message: %v", err)
		}
	}
}
