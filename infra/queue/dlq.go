package queue

import (
	"context"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (q *Queue) StartDLQProcessor(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	dlqName := q.queueName + "_dlq"

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping DLQ processor")
			return
		case <-ticker.C:
			ch, _ := q.getChannel()
			if ch == nil || ch.IsClosed() {
				continue
			}

			for {
				msg, ok, err := ch.Get(dlqName, false)
				if err != nil {
					slog.Error("Error fetching from DLQ", "error", err)
					break
				}
				if !ok {
					break
				}

				var retryCount int32
				if msg.Headers != nil {
					if count, ok := msg.Headers["x-retry-count"].(int32); ok {
						retryCount = count
					}
				}

				if retryCount >= 5 {
					bodyLimit := 100
					if len(msg.Body) < bodyLimit {
						bodyLimit = len(msg.Body)
					}
					slog.Warn("Message exceeded max retries, dropping permanently", "retry_count", retryCount, "body_snippet", string(msg.Body[:bodyLimit]))
					msg.Ack(false)
					continue
				}

				headers := msg.Headers
				if headers == nil {
					headers = make(amqp.Table)
				}
				headers["x-retry-count"] = retryCount + 1

				err = ch.Publish(
					"",
					q.queueName,
					false,
					false,
					amqp.Publishing{
						Headers:      headers,
						ContentType:  msg.ContentType,
						Body:         msg.Body,
						DeliveryMode: msg.DeliveryMode,
					},
				)

				if err != nil {
					slog.Error("Error requeuing message from DLQ", "error", err)
					msg.Nack(false, true)
					break
				} else {
					slog.Info("Successfully requeued a message from DLQ")
					msg.Ack(false)
				}
			}
		}
	}
}
