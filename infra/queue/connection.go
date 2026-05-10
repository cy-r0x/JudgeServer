package queue

import (
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (q *Queue) connect() error {
	var err error
	conn, err := amqp.Dial(q.rabbitmqURL)
	if err != nil {
		slog.Error("Failed to connect to RabbitMQ", "error", err)
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		slog.Error("Failed to open channel", "error", err)
		conn.Close()
		return err
	}

	dlxName := q.queueName + "_dlx"
	err = ch.ExchangeDeclare(dlxName, "direct", true, false, false, false, nil)
	if err != nil {
		slog.Error("Failed to declare DLX", "error", err)
		ch.Close()
		conn.Close()
		return err
	}

	dlqName := q.queueName + "_dlq"
	_, err = ch.QueueDeclare(dlqName, true, false, false, false, amqp.Table{"x-queue-type": "quorum"})
	if err != nil {
		slog.Error("Failed to declare DLQ", "error", err)
		ch.Close()
		conn.Close()
		return err
	}

	err = ch.QueueBind(dlqName, q.queueName, dlxName, false, nil)
	if err != nil {
		slog.Error("Failed to bind DLQ to DLX", "error", err)
		ch.Close()
		conn.Close()
		return err
	}

	args := amqp.Table{
		"x-queue-type":              "quorum",
		"x-dead-letter-exchange":    dlxName,
		"x-dead-letter-routing-key": q.queueName,
	}
	_, err = ch.QueueDeclare(q.queueName, true, false, false, false, args)
	if err != nil {
		slog.Error("Failed to declare queue", "error", err)
		ch.Close()
		conn.Close()
		return err
	}

	q.mu.Lock()
	if q.ch != nil {
		q.ch.Close()
	}
	if q.conn != nil {
		q.conn.Close()
	}
	q.conn = conn
	q.ch = ch
	q.mu.Unlock()

	return nil
}

func (q *Queue) reconnect() error {
	slog.Info("Attempting to reconnect to RabbitMQ")

	if q.ch != nil {
		q.ch.Close()
	}
	if q.conn != nil {
		q.conn.Close()
	}

	backoff := time.Second
	maxBackoff := 30 * time.Second

	for {
		if q.ctx != nil {
			select {
			case <-q.ctx.Done():
				return q.ctx.Err()
			default:
			}
		}

		err := q.connect()
		if err == nil {
			slog.Info("Successfully reconnected to RabbitMQ")
			return nil
		}

		slog.Error("Reconnection failed", "backoff", backoff, "error", err)

		if q.ctx != nil {
			timer := time.NewTimer(backoff)
			select {
			case <-timer.C:
			case <-q.ctx.Done():
				timer.Stop()
				return q.ctx.Err()
			}
		} else {
			time.Sleep(backoff)
		}

		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
		}
	}
}
