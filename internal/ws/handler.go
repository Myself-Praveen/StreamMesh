package ws

import (
	"encoding/json"

	"github.com/Myself-Praveen/StreamMesh/internal/logger"
	"go.uber.org/zap"
)

// MessageHandler handles incoming WebSocket protocol messages
type MessageHandler struct {
	manager  *Manager
	registry pubsubRegistry
	router   pubsubRouter
}

// Interfaces to break circular dependencies
type pubsubRegistry interface {
	Subscribe(topicName, connID string)
	Unsubscribe(topicName, connID string)
}

type pubsubRouter interface {
	Publish(topic string, payload []byte) error
}

// NewMessageHandler creates a new handler
func NewMessageHandler(manager *Manager, registry pubsubRegistry, router pubsubRouter) *MessageHandler {
	return &MessageHandler{
		manager:  manager,
		registry: registry,
		router:   router,
	}
}

// HandleMessage parses and processes a WebSocket message
func (h *MessageHandler) HandleMessage(conn *Connection, data []byte) {
	env, err := ParseMessage(data)
	if err != nil {
		logger.Log.Warn("Invalid message format", zap.Error(err), zap.String("id", conn.ID))
		return
	}

	switch env.Type {
	case TypeSubscribe:
		h.registry.Subscribe(env.Topic, conn.ID)
		conn.AddTopic(env.Topic)
		logger.Log.Debug("Client subscribed", zap.String("id", conn.ID), zap.String("topic", env.Topic))

	case TypeUnsubscribe:
		h.registry.Unsubscribe(env.Topic, conn.ID)
		conn.RemoveTopic(env.Topic)
		logger.Log.Debug("Client unsubscribed", zap.String("id", conn.ID), zap.String("topic", env.Topic))

	case TypePublish:
		// Payload is raw JSON in the envelope
		if err := h.router.Publish(env.Topic, env.Payload); err != nil {
			logger.Log.Error("Failed to publish message", zap.Error(err), zap.String("topic", env.Topic))
		}

	case TypeHistory:
		// TODO: handle getting history
		logger.Log.Debug("History requested", zap.String("id", conn.ID), zap.String("topic", env.Topic))

	default:
		logger.Log.Warn("Unknown message type", zap.String("type", string(env.Type)))
	}
}
