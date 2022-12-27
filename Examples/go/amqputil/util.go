package amqputil

import (
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PrintMessage(delivery *amqp.Delivery, includeProperties bool) {
	fmt.Printf("Exchange: %s\n", delivery.Exchange)
	fmt.Printf("Routing key: %s\n", delivery.RoutingKey)

	if includeProperties {
		if delivery.ContentType != "" {
			fmt.Printf("Content type: %s\n", delivery.ContentType)
		}

		if delivery.ContentEncoding != "" {
			fmt.Printf("Content encoding: %s\n", delivery.ContentEncoding)
		}

		if delivery.DeliveryMode != 0 {
			fmt.Printf("Delivery mode: %d\n", delivery.DeliveryMode)
		}

		if delivery.Priority != 0 {
			fmt.Printf("Priority: %d\n", delivery.Priority)
		}

		if delivery.CorrelationId != "" {
			fmt.Printf("Correlation id: %s\n", delivery.CorrelationId)
		}

		if delivery.ReplyTo != "" {
			fmt.Printf("Reply to: %s\n", delivery.ReplyTo)
		}

		if delivery.Expiration != "" {
			fmt.Printf("Expiration: %s\n", delivery.Expiration)
		}

		if delivery.MessageId != "" {
			fmt.Printf("Message id: %s\n", delivery.MessageId)
		}

		if !delivery.Timestamp.Equal(time.Time{}) {
			fmt.Printf("Timestamp: %s\n", delivery.Timestamp)
		}

		if delivery.Type != "" {
			fmt.Printf("Type: %s\n", delivery.Type)
		}

		if delivery.UserId != "" {
			fmt.Printf("User id: %s\n", delivery.UserId)
		}

		if delivery.AppId != "" {
			fmt.Printf("App id: %s\n", delivery.AppId)
		}

		if len(delivery.Headers) != 0 {
			b, err := json.MarshalIndent(delivery.Headers, "", " ")
			if err != nil {
				fmt.Printf("Failed to marshal message headers: %v", err)
			}

			fmt.Printf("Headers: %s\n", b)
		}
	}

	fmt.Printf("\nBody: %s\n", delivery.Body)
}
