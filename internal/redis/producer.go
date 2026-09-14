package redis

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

// StreamMessage represents a message to be published to a Redis stream
type StreamMessage struct {
	Topic      string `json:"topic"`
	Payload    []byte `json:"payload"`
	NodeID     string `json:"node_id"`
	EnvelopeID string `json:"envelope_id"`
}

// ProduceMessage publishes a message to a Redis Stream
func ProduceMessage(ctx context.Context, stream string, msg StreamMessage) error {
	if Client == nil {
		return nil // NOOP if Redis is not initialized
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	args := &redis.XAddArgs{
		Stream: stream,
		Values: map[string]interface{}{
			"data": data,
		},
	}

	return Client.XAdd(ctx, args).Err()
}
