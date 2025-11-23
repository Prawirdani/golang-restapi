package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prawirdani/golang-restapi/internal/domain/auth"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	AuthDirectExchange           = "auth.direct"
	ResetPasswordEmailRoutingKey = "email.reset-password"
	ResetPasswordEmailQueue      = "auth.email.reset-password"
)

var ResetPasswordEmailTopology = &Topology{
	Name:         "Reset Password Email Topology",
	Exchange:     AuthDirectExchange,
	ExchangeType: "direct",
	Queue:        ResetPasswordEmailQueue,
	RoutingKey:   ResetPasswordEmailRoutingKey,
	Durable:      true,
	RetryTTL:     5000, // 5 Seconds
	MaxRetry:     3,
	QueueArgs: amqp.Table{
		"x-queue-type": "quorum",
	},
}

type AuthMessagePublisher struct {
	conn *amqp.Connection
}

func NewAuthMessagePublisher(conn *amqp.Connection) *AuthMessagePublisher {
	return &AuthMessagePublisher{conn: conn}
}

// Implements auth.MessagePublisher
func (mp *AuthMessagePublisher) SendResetPasswordEmail(
	ctx context.Context,
	msg auth.ResetPasswordEmailMessage,
) error {
	// NOTE: For low to moderate traffic is okay to open channel per function call, but when the traffic goes up it
	// slightly more overhead per publish (channel open/close is a network round-trip)
	// TODO: Use thread safe channel or use channel pool
	ch, err := mp.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(
		ctx,
		AuthDirectExchange,
		ResetPasswordEmailRoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
			Timestamp:   time.Now(),
			MessageId:   uuid.NewString(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish reset password email message: %w", err)
	}

	return nil
}
