package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/prawirdani/golang-restapi/internal/infra/mq"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RMQPublisher can publish any payload to any queue
type RMQPublisher struct {
	conn *amqp.Connection
}

func NewPublisher(url string) (*RMQPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &RMQPublisher{
		conn: conn,
		// channel: ch,
	}, nil
}

func (p *RMQPublisher) Publish(
	ctx context.Context,
	queueName string,
	payload any,
) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	// Declare queue (idempotent)
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Marshal payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create message
	msg := mq.NewMessage(payloadBytes)

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Publish
	err = ch.PublishWithContext(
		ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}

	log.Printf("Published [Queue: %s,  ID: %s]", queueName, msg.ID)
	return nil
}

func (p *RMQPublisher) Close() error {
	if p.conn == nil {
		return nil
	}
	return p.conn.Close()
}
