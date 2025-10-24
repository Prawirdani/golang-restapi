package mq

import "context"

type MessageProducer interface {
	Publish(ctx context.Context, queueName string, payload any) error
	Close() error
}

type MessageConsumer interface {
	RegisterHandler(queueName string, handler MessageHandler)
	Start(ctx context.Context) error
	Close() error
}
