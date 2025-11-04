package consumer

import (
	"context"
	"errors"
	"fmt"

	"github.com/prawirdani/golang-restapi/pkg/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type HandlerFunc func(amqp.Delivery) error

// TODO: Retry Mechanism
type Registry struct {
	conn      *amqp.Connection
	consumers map[string]HandlerFunc // Queue Name -> Handler/Consumer func
}

func NewRegistry(conn *amqp.Connection) *Registry {
	return &Registry{conn: conn}
}

func (r *Registry) RegisterConsumers(consumers map[string]HandlerFunc) {
	r.consumers = consumers
}

func (r *Registry) Start(ctx context.Context) error {
	if r.consumers == nil {
		return errors.New("consumers is nil")
	}

	errCh := make(chan error, len(r.consumers)) // buffer for all handlers

	for queue, handler := range r.consumers {
		go func() {
			if err := r.consume(ctx, queue, handler); err != nil {
				errCh <- err
			}
		}()
	}

	// Wait for first error or context cancellation
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// consume is the generic consumer that handles the boilerplate
func (r *Registry) consume(
	ctx context.Context,
	queueName string,
	handler HandlerFunc,
) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}
	defer ch.Close()

	// Set QoS
	if err := ch.Qos(1, 0, false); err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	msgs, err := ch.ConsumeWithContext(
		ctx,
		queueName,
		"",    // consumer tag
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to consume %s queue: %w", queueName, err)
	}

	log.Info("Consumer started", "queue", queueName)

	for {
		select {
		case <-ctx.Done():
			log.Info("Consumer shutting down", "queue", queueName)
			return ctx.Err()

		case d, ok := <-msgs:
			if !ok {
				return fmt.Errorf("message channel closed")
			}

			log.Info("Message Received", "queue", queueName)

			if err := handler(d); err != nil {
				log.ErrorCtx(ctx, "Failed to process message",
					"error", err.Error(),
					"queue", queueName,
				)
				d.Nack(false, true) // requeue on error
			} else {
				log.Info("Message Processed", "queue", queueName)
				d.Ack(false)
			}
		}
	}
}
