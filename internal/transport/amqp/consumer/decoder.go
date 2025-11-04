package consumer

import (
	"encoding/json"
	"fmt"
)

func decodeJsonBody[T any](b []byte) (*T, error) {
	var body T
	err := json.Unmarshal(b, &body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode message body: %w", err)
	}

	return &body, nil
}
