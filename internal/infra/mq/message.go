package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type MessageHandler func(ctx context.Context, payload json.RawMessage) error

// Message wraps any payload with metadata
type Message struct {
	ID        string          `json:"id"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}

func NewMessage(payload json.RawMessage) Message {
	return Message{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Timestamp: time.Now(),
		Payload:   payload,
	}
}
