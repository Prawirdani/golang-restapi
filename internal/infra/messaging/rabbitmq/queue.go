package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

func (queue *Queue) Declare(ch *amqp.Channel) error {
	_, err := ch.QueueDeclare(
		queue.Name,
		queue.Durable,
		queue.AutoDelete,
		queue.Exclusive,
		queue.NoWait,
		queue.Args,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", queue.Name, err)
	}

	return nil
}
