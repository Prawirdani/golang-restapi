package auth

import (
	"context"

	"github.com/prawirdani/golang-restapi/internal/messages"
)

// Interface for message/event producer, eg rabbitmq, kafka etc.
type MessagePublisher interface {
	SendResetPasswordEmail(ctx context.Context, msg messages.ResetPasswordEmail) error
}
