package consumer

import (
	"context"
	"fmt"
	"maps"
	"time"

	"github.com/prawirdani/golang-restapi/internal/infra/messaging/rabbitmq"
	"github.com/prawirdani/golang-restapi/pkg/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type HandlerFunc func(ctx context.Context, d amqp.Delivery) error

type ConsumerClient struct {
	conn *amqp.Connection
}

func NewConsumerClient(conn *amqp.Connection) *ConsumerClient {
	return &ConsumerClient{conn: conn}
}

func (m *ConsumerClient) Consume(
	ctx context.Context,
	tpl *rabbitmq.Topology,
	handler HandlerFunc,
) error {
	ctx = log.WithContext(ctx, "queue", tpl.Queue)

	ch, err := m.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	if err := ch.Qos(1, 0, false); err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	msgs, err := ch.ConsumeWithContext(ctx, tpl.Queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to consume from %s: %w", tpl.Queue, err)
	}

	log.InfoCtx(ctx, "Consumer started")

	for {
		select {
		case <-ctx.Done():
			log.InfoCtx(ctx, "Consumer shutting down")
			return ctx.Err()

		case d, ok := <-msgs:
			if !ok {
				return fmt.Errorf("message channel closed: %s", tpl.Queue)
			}
			ctx := log.WithContext(ctx, "message_id", d.MessageId)

			log.InfoCtx(ctx, "Message received")
			if err := handler(ctx, d); err != nil {
				log.ErrorCtx(ctx, "Failed to handle message", err)
				m.handleFailure(ctx, ch, tpl, d, err)
			} else {
				if err := d.Ack(false); err != nil {
					log.ErrorCtx(ctx, "Failed to Ack message", err)
				}
				log.InfoCtx(ctx, "Message processed")
			}
		}
	}
}

func (m *ConsumerClient) handleFailure(
	ctx context.Context,
	ch *amqp.Channel,
	tpl *rabbitmq.Topology,
	d amqp.Delivery,
	err error,
) {
	deathCount := m.getDeathCount(d.Headers)

	// If max retry atempts, send the message to dlq
	if deathCount >= int64(tpl.MaxRetry) {
		log.InfoCtx(ctx, "Max retry attempts, message sent to DLQ")
		headers := amqp.Table{
			"x-error":          err.Error(),      // inject error to headers to inspect latter
			"x-original-queue": tpl.Queue,        // track source queue
			"x-failed-at":      time.Now().UTC(), // timestamp of failure
		}
		maps.Copy(headers, d.Headers)
		_ = ch.PublishWithContext(
			ctx,
			tpl.Exchange+rabbitmq.DLXSuffix,
			tpl.RoutingKey,
			false,
			false,
			amqp.Publishing{
				Headers:       headers,
				ContentType:   d.ContentType,
				Body:          d.Body,
				DeliveryMode:  d.DeliveryMode,
				CorrelationId: d.CorrelationId,
			},
		)
		// Ack Message
		if err := d.Ack(false); err != nil {
			log.ErrorCtx(ctx, "Failed to Ack message", err)
		}
		return
	}
	log.InfoCtx(ctx, "Message sent to retry queue")
	if err := d.Nack(false, false); err != nil {
		log.ErrorCtx(ctx, "Failed to Nack message", err)
	}
}

func (m *ConsumerClient) getDeathCount(table amqp.Table) int64 {
	deaths, ok := table["x-death"].([]any)
	if !ok || len(deaths) == 0 {
		return 0
	}

	first, ok := deaths[0].(amqp.Table)
	if !ok {
		return 0
	}

	count, _ := first["count"].(int64)
	return count
}
