package queue

import (
	"context"
	"sync"

	"github.com/judgenot0/judge-backend/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue struct {
	msgs        <-chan amqp.Delivery
	conn        *amqp.Connection
	ch          *amqp.Channel
	queueName   string
	rabbitmqURL string
	workerCount int
	ctx         context.Context
	mu          sync.RWMutex
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) InitQueue(config *config.Config) error {
	q.queueName = config.QueueName
	q.rabbitmqURL = config.RabbitMQURL

	return q.connect()
}

func (q *Queue) getChannel() (*amqp.Channel, *amqp.Connection) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.ch, q.conn
}

func (q *Queue) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()
	var errs []error
	if q.ch != nil {
		if err := q.ch.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if q.conn != nil {
		if err := q.conn.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
