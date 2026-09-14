package ws

import (
	"sync"
	"time"

	"github.com/Myself-Praveen/StreamMesh/internal/config"
	"github.com/Myself-Praveen/StreamMesh/internal/ratelimit"
	"github.com/gorilla/websocket"
)

// Connection represents a single WebSocket client connection
type Connection struct {
	ID            string
	Conn          *websocket.Conn
	Send          chan []byte
	Topics        map[string]bool
	CreatedAt     time.Time
	LastPing      time.Time
	Metadata      map[string]string
	Subscriptions map[string]bool
	mu            sync.RWMutex
	limiter       *ratelimit.TokenBucket
}

// NewConnection creates a new Connection instance
func NewConnection(id string, conn *websocket.Conn) *Connection {
	var limiter *ratelimit.TokenBucket
	if config.AppConfig != nil && config.AppConfig.WebSocket.RateLimit.Burst > 0 {
		limiter = ratelimit.NewTokenBucket(
			config.AppConfig.WebSocket.RateLimit.Burst,
			float64(config.AppConfig.WebSocket.RateLimit.MessagesPerSecond),
		)
	} else {
		// Default limits if config not loaded
		limiter = ratelimit.NewTokenBucket(100, 50.0)
	}

	return &Connection{
		ID:            id,
		Conn:          conn,
		Send:          make(chan []byte, 256),
		Topics:        make(map[string]bool),
		CreatedAt:     time.Now(),
		LastPing:      time.Now(),
		Metadata:      make(map[string]string),
		Subscriptions: make(map[string]bool),
		limiter:       limiter,
	}
}

// AddTopic subscribes the connection to a topic
func (c *Connection) AddTopic(topic string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Topics[topic] = true
}

// RemoveTopic unsubscribes the connection from a topic
func (c *Connection) RemoveTopic(topic string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.Topics, topic)
}

// GetTopics returns a list of topics the connection is subscribed to
func (c *Connection) GetTopics() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	topics := make([]string, 0, len(c.Topics))
	for t := range c.Topics {
		topics = append(topics, t)
	}
	return topics
}

// UpdateLastPing updates the last ping time to now
func (c *Connection) UpdateLastPing() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.LastPing = time.Now()
}

// GetLastPing returns the last ping time
func (c *Connection) GetLastPing() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.LastPing
}

// Allow checks the token bucket rate limiter
func (c *Connection) Allow() bool {
	if c.limiter != nil {
		return c.limiter.Allow()
	}
	return true
}

// SetMetadata sets a metadata key-value pair
func (c *Connection) SetMetadata(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Metadata[key] = value
}

// GetMetadata retrieves a metadata value by key
func (c *Connection) GetMetadata(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.Metadata[key]
	return val, ok
}
