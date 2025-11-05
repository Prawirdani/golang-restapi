package main

import (
	"context"
	"fmt"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/infra/messaging/rabbitmq"
	"github.com/prawirdani/golang-restapi/internal/transport/amqp/consumer"
	"github.com/prawirdani/golang-restapi/pkg/mailer"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Start message consumers, this function is blocking, so run it inside seperate goroutine
func startMessageConsumers(
	ctx context.Context,
	conn *amqp.Connection,
	cfg *config.Config,
) error {
	m := mailer.New(cfg)
	authConsumers := consumer.NewAuthMessageConsumer(m)
	consumerClient := consumer.NewConsumerClient(conn)

	errCh := make(chan error, 1)

	// Run consumers in background
	// TODO: As things grows, consider using slice of [topology+handler] and run all of it through loops
	go func() {
		if err := consumerClient.Consume(
			ctx,
			rabbitmq.ResetPasswordEmailTopology,
			authConsumers.EmailResetPasswordHandler,
		); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()

	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("consumer error: %w", err)
		}
		return nil
	}
}
