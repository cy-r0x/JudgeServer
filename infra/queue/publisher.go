package queue

import (
	"log/slog"

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
		slog.Error("Failed to publish message, attempting reconnect", "error", err)
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
