package pubsub

import (
	"context"
	"fmt"
	"time"

	"github.com/Myself-Praveen/StreamMesh/internal/logger"
	"github.com/Myself-Praveen/StreamMesh/internal/redis"
	"github.com/Myself-Praveen/StreamMesh/internal/ws"
	"go.uber.org/zap"
)

// Router handles publishing messages to subscribers
type Router struct {
	registry *TopicRegistry
	manager  *ws.Manager
}

// ClusterPublisher interface for broadcasting messages across nodes
type ClusterPublisher interface {
	BroadcastMessage(topic string, payload []byte, envelopeID string)
}

// GlobalClusterPublisher is an optional gRPC publisher
var GlobalClusterPublisher ClusterPublisher

// GlobalRouter is the default message router
var GlobalRouter *Router

// NewRouter creates a new local pub/sub router
func NewRouter(registry *TopicRegistry, manager *ws.Manager) *Router {
	router := &Router{
		registry: registry,
		manager:  manager,
	}
	GlobalRouter = router
	return router
}

// PublishLocal routes a message to all subscribers of a topic locally
func (r *Router) PublishLocal(topic string, payload []byte) error {
	subs := r.registry.GetSubscribers(topic)
	if len(subs) == 0 {
		return nil
	}

	// Create the envelope
	envID := fmt.Sprintf("msg-%d", time.Now().UnixNano())
	env, err := ws.NewMessage(ws.TypePublish, topic, nil, envID)
	if err != nil {
		return err
	}
	
	// Set the pre-marshaled payload since we might receive raw bytes
	env.Payload = payload

	b, err := env.Encode()
	if err != nil {
		return err
	}

	// Add to topic history
	t := r.registry.getOrCreateTopic(topic)
	t.History.Add(b)

	delivered := 0
	for _, connID := range subs {
		conn, ok := r.manager.Get(connID)
		if !ok {
			// Subscriber disconnected, remove from registry
			r.registry.Unsubscribe(topic, connID)
			continue
		}

		select {
		case conn.Send <- b:
			delivered++
		default:
			// Buffer full, drop message or disconnect slow client
			logger.Log.Warn("Send buffer full, dropping message", zap.String("id", connID))
		}
	}

	logger.Log.Debug("Published message locally", zap.String("topic", topic), zap.Int("delivered", delivered))
	return nil
}

// Publish routes a message to a topic (local and distributed)
func (r *Router) Publish(topic string, payload []byte) error {
	// Publish locally first for lower latency
	err := r.PublishLocal(topic, payload)

	// Publish to gRPC Publisher if available
	envID := fmt.Sprintf("env-%d", time.Now().UnixNano())
	if GlobalClusterPublisher != nil {
		GlobalClusterPublisher.BroadcastMessage(topic, payload, envID)
	} else if redis.Client != nil {
		// Fallback to Redis stream for fan-out
		streamMsg := redis.StreamMessage{
			Topic:      topic,
			Payload:    payload,
			NodeID:     redis.CurrentNode.ID,
			EnvelopeID: envID,
		}
		// Using background context since this is asynchronous fire-and-forget
		redisErr := redis.ProduceMessage(context.Background(), "streammesh:messages", streamMsg)
		if redisErr != nil {
			logger.Log.Error("Failed to publish to Redis stream", zap.Error(redisErr))
		}
	}

	return err
}
