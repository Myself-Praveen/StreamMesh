package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Myself-Praveen/StreamMesh/internal/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Consumer handles reading from a Redis stream via a consumer group
type Consumer struct {
	stream       string
	group        string
	consumerName string
	handler      func(StreamMessage)
}

// NewConsumer creates a new Redis stream consumer
func NewConsumer(stream, group, consumerName string, handler func(StreamMessage)) *Consumer {
	return &Consumer{
		stream:       stream,
		group:        group,
		consumerName: consumerName,
		handler:      handler,
	}
}

// Start starts consuming messages in a separate goroutine
func (c *Consumer) Start(ctx context.Context) {
	if Client == nil {
		return
	}

	// Ensure group exists
	err := Client.XGroupCreateMkStream(ctx, c.stream, c.group, "$").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		logger.Log.Error("Failed to create consumer group", zap.Error(err))
		return
	}

	go c.consumeLoop(ctx)
}

func (c *Consumer) consumeLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			args := &redis.XReadGroupArgs{
				Group:    c.group,
				Consumer: c.consumerName,
				Streams:  []string{c.stream, ">"},
				Count:    10,
				Block:    2 * time.Second,
			}

			streams, err := Client.XReadGroup(ctx, args).Result()
			if err != nil {
				if err != redis.Nil { // redis.Nil is returned on timeout
					logger.Log.Error("Error reading from stream", zap.Error(err))
					time.Sleep(1 * time.Second) // backoff
				}
				continue
			}

			for _, stream := range streams {
				for _, msg := range stream.Messages {
					dataStr, ok := msg.Values["data"].(string)
					if !ok {
						logger.Log.Warn("Invalid data format in stream message")
						continue
					}

					var streamMsg StreamMessage
					if err := json.Unmarshal([]byte(dataStr), &streamMsg); err != nil {
						logger.Log.Error("Failed to unmarshal stream message", zap.Error(err))
						continue
					}

					// Process message
					if c.handler != nil {
						c.handler(streamMsg)
					}

					// Acknowledge message
					Client.XAck(ctx, c.stream, c.group, msg.ID)
				}
			}
		}
	}
}
