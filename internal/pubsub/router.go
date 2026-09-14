package pubsub

import (
	"fmt"
	"time"

	"github.com/Myself-Praveen/StreamMesh/internal/logger"
	"github.com/Myself-Praveen/StreamMesh/internal/ws"
	"go.uber.org/zap"
)

// Router handles publishing messages to subscribers
type Router struct {
	registry *TopicRegistry
	manager  *ws.Manager
}

// NewRouter creates a new local pub/sub router
func NewRouter(registry *TopicRegistry, manager *ws.Manager) *Router {
	return &Router{
		registry: registry,
		manager:  manager,
	}
}

// Publish routes a message to all subscribers of a topic locally
func (r *Router) Publish(topic string, payload []byte) error {
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

	logger.Log.Debug("Published message", zap.String("topic", topic), zap.Int("delivered", delivered))
	return nil
}
