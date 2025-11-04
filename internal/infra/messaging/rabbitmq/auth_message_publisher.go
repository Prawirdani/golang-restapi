package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/prawirdani/golang-restapi/internal/messages"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AuthMessagePublisher struct {
	ch *amqp.Channel
}

func NewAuthMessagePublisher(conn *amqp.Connection) *AuthMessagePublisher {
	ch, err := conn.Channel()
	if err != nil {
		return nil
	}

	return &AuthMessagePublisher{ch: ch}
}

// Implements auth.MessagePublisher
func (mp *AuthMessagePublisher) SendResetPasswordEmail(
	ctx context.Context,
	msg messages.ResetPasswordEmail,
) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = mp.ch.PublishWithContext(
		ctx,
		AuthExchange.Name,
		ResetPasswordEmailRoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
			Timestamp:   time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish reset password email message: %w", err)
	}

	return nil
}
