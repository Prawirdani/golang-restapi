package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Exchange struct {
	Name       string
	Type       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

func (exc *Exchange) Declare(ch *amqp.Channel) error {
	if err := ch.ExchangeDeclare(
		exc.Name,
		exc.Type,
		exc.Durable,
		exc.AutoDelete,
		exc.Internal,
		exc.NoWait,
		exc.Args,
	); err != nil {
		return fmt.Errorf("failed to declare exchange %s: %w", exc.Name, err)
	}

	return nil
}

func (exc *Exchange) BindQueue(ch *amqp.Channel, queue Queue, routingKey string) error {
	err := ch.QueueBind(queue.Name, routingKey, exc.Name, false, nil)
	if err != nil {
		return fmt.Errorf("failed to bind exchange %s with %s queue: %w", exc.Name, queue.Name, err)
	}
	return nil
}
