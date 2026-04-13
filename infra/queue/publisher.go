package queue

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (q *Queue) QueueMessage(submission []byte) error {
	ch, _ := q.getChannel()
	if ch == nil || ch.IsClosed() {
		if err := q.reconnect(); err != nil {
			return err
		}
		ch, _ = q.getChannel()
	}

	err := ch.Publish(
		"",
		q.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        submission,
		},
	)

	if err != nil {
		log.Printf("Failed to publish message, attempting reconnect: %v", err)
		if reconnectErr := q.reconnect(); reconnectErr != nil {
			return reconnectErr
		}
		ch, _ = q.getChannel()
		err = ch.Publish(
			"",
			q.queueName,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        submission,
			},
		)
	}

	return err
}
