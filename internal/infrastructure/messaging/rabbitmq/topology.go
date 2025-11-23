package rabbitmq

import (
	"fmt"
	"maps"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Topology defines the complete RabbitMQ topology for a given message type.
// It includes the main exchange/queue, a retry exchange/queue, and a dead-letter exchange/queue (DLX/DLQ).
//
// Every Topology implicitly follows this flow:
//
//	main queue → retry queue (via x-dead-letter-exchange) → back to main queue → eventually DLQ after MaxRetry.
//
// Example flow:
//  1. Message is published to the main exchange → main queue.
//  2. Consumer fails and NACKs message → automatically sent to retry exchange → retry queue.
//  3. Message sits in retry queue for RetryTTL → returns to main queue.
//  4. Repeat until x-death count reaches MaxRetry → manually published to DLQ.
type Topology struct {
	Name         string
	Exchange     string
	ExchangeType string

	Queue      string
	QueueArgs  amqp.Table
	RoutingKey string

	// Common queue/exchange properties
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool

	// Retry settings
	RetryTTL int32 // Duration (ms) messages are held in retry queue before re-queued
	MaxRetry int   // Max retry attempts based on x-death header
}

// Declare creates the full topology in RabbitMQ: exchanges, queues, and bindings.
// It declares the main exchange/queue, the retry exchange/queue, and the DLX/DLQ.
//
// Notes:
// - The retry mechanism uses x-dead-letter-exchange and a TTL on the retry queue.
// - DLQ messages are published manually after MaxRetry attempts, using the x-death header as the counter.
// - All queues are durable and have their properties set according to the Topology fields.
// - Bindings connect queues to their respective exchanges with the RoutingKey.
func (t *Topology) Declare(ch *amqp.Channel) error {
	mainExchange := t.Exchange
	mainQueue := t.Queue
	retryExchange := mainExchange + RetrySuffix
	dlxExchange := mainExchange + DLXSuffix

	retryQueue := mainQueue + RetrySuffix
	dlqQueue := mainQueue + DLQSuffix

	// 1. Declare main exchange
	if err := ch.ExchangeDeclare(
		mainExchange,
		t.ExchangeType,
		t.Durable,
		t.AutoDelete,
		false,
		t.NoWait,
		nil,
	); err != nil {
		return fmt.Errorf("declare main exchange: %w", err)
	}

	// 2. Declare retry exchange
	if err := ch.ExchangeDeclare(
		retryExchange,
		t.ExchangeType,
		t.Durable,
		t.AutoDelete,
		false,
		t.NoWait,
		nil,
	); err != nil {
		return fmt.Errorf("declare retry exchange: %w", err)
	}

	// 3. Declare dead letter exchange
	if err := ch.ExchangeDeclare(
		dlxExchange,
		t.ExchangeType,
		t.Durable,
		t.AutoDelete,
		false,
		t.NoWait,
		nil,
	); err != nil {
		return fmt.Errorf("declare dlx exchange: %w", err)
	}

	mainArgs := amqp.Table{}
	maps.Copy(mainArgs, t.QueueArgs)

	mainArgs["x-dead-letter-exchange"] = retryExchange
	mainArgs["x-dead-letter-routing-key"] = t.RoutingKey

	if _, err := ch.QueueDeclare(
		mainQueue,
		t.Durable,
		t.AutoDelete,
		t.Exclusive,
		t.NoWait,
		mainArgs,
	); err != nil {
		return fmt.Errorf("declare main queue: %w", err)
	}

	retryArgs := amqp.Table{}
	maps.Copy(retryArgs, t.QueueArgs)
	retryArgs["x-dead-letter-exchange"] = mainExchange
	retryArgs["x-dead-letter-routing-key"] = t.RoutingKey
	retryArgs["x-message-ttl"] = t.RetryTTL
	if _, err := ch.QueueDeclare(
		retryQueue,
		t.Durable,
		t.AutoDelete,
		t.Exclusive,
		t.NoWait,
		retryArgs,
	); err != nil {
		return fmt.Errorf("declare retry queue: %w", err)
	}

	dlqArgs := amqp.Table{}
	if queueType, exists := t.QueueArgs["x-queue-type"]; exists {
		dlqArgs["x-queue-type"] = queueType
	}

	if _, err := ch.QueueDeclare(
		dlqQueue,
		t.Durable,
		t.AutoDelete,
		t.Exclusive,
		t.NoWait,
		dlqArgs,
	); err != nil {
		return fmt.Errorf("declare DLQ: %w", err)
	}

	// 5. Bind queues
	if err := ch.QueueBind(mainQueue, t.RoutingKey, mainExchange, false, nil); err != nil {
		return fmt.Errorf("bind main queue: %w", err)
	}

	if err := ch.QueueBind(retryQueue, t.RoutingKey, retryExchange, false, nil); err != nil {
		return fmt.Errorf("bind retry queue: %w", err)
	}

	if err := ch.QueueBind(dlqQueue, t.RoutingKey, dlxExchange, false, nil); err != nil {
		return fmt.Errorf("bind DLQ: %w", err)
	}

	return nil
}

func SetupTopologies(conn *amqp.Connection, topologies ...*Topology) error {
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("open channel: %w", err)
	}
	defer ch.Close()

	for _, t := range topologies {
		if err := t.Declare(ch); err != nil {
			return fmt.Errorf("declare topology %s: %w", t.Name, err)
		}
	}
	return nil
}
