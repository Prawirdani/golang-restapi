package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/prawirdani/golang-restapi/internal/infra/mq"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer can handle multiple job types from multiple queues
type Consumer struct {
	conn     *amqp.Connection
	handlers map[string]mq.MessageHandler // queueName:handler
	mu       sync.RWMutex
}

func NewConsumer(url string) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &Consumer{
		conn:     conn,
		handlers: make(map[string]mq.MessageHandler),
	}, nil
}

func (c *Consumer) RegisterHandler(queueName string, handler mq.MessageHandler) {
	c.mu.Lock()
	c.handlers[queueName] = handler
	c.mu.Unlock()

	log.Printf("Registered handler for: %s", queueName)
}

// Start begins consuming from all registered queues
func (c *Consumer) Start(ctx context.Context) error {
	var wg sync.WaitGroup

	// Each message type gets its own queue + goroutine
	for msgType := range c.handlers {
		wg.Add(1)
		go func(queueName string) {
			defer wg.Done()
			c.consumeQueue(ctx, queueName) // queueName = msgType
		}(msgType)
	}

	wg.Wait()
	return nil
}

func (c *Consumer) consumeQueue(ctx context.Context, queueName string) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	// Declare queue
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Set QoS
	err = ch.Qos(1, 0, false)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to consume: %w", err)
	}

	log.Printf("Consumer started for queue: %s", queueName)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case delivery, ok := <-msgs:
			if !ok {
				return fmt.Errorf("channel closed")
			}

			c.handleMessage(ctx, queueName, delivery)
		}
	}
}

func (c *Consumer) handleMessage(ctx context.Context, queueName string, delivery amqp.Delivery) {
	// Unmarshal message
	var msg mq.Message
	err := json.Unmarshal(delivery.Body, &msg)
	if err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		delivery.Nack(false, false)
		return
	}

	log.Printf("Received [Queue: %s, ID: %s]", queueName, msg.ID)

	// Find handler
	c.mu.RLock()
	handler, exists := c.handlers[queueName]
	c.mu.RUnlock()

	if !exists {
		log.Printf("No handler for queue: %s", queueName)
		delivery.Nack(false, false) // Don't requeue unknown types
		return
	}

	// Execute handler
	err = handler(ctx, msg.Payload)
	if err != nil {
		log.Printf("Handler error: %v", err)
		delivery.Nack(false, true) // Requeue on error
	} else {
		delivery.Ack(false)
		log.Printf("Message processed [Queue: %s, ID: %s]", queueName, msg.ID)
	}
}

func (c *Consumer) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}
