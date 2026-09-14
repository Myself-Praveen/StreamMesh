package pubsub

import (
	"sync"
)

// TopicRegistry manages all active topics and their subscribers
type TopicRegistry struct {
	topics map[string]*Topic
	mu     sync.RWMutex
}

// Topic represents a single pub/sub topic
type Topic struct {
	Name        string
	Subscribers map[string]bool // map of Connection IDs
	History     *HistoryBuffer
	mu          sync.RWMutex
}

// NewTopicRegistry creates a new registry
func NewTopicRegistry() *TopicRegistry {
	return &TopicRegistry{
		topics: make(map[string]*Topic),
	}
}

// getOrCreateTopic returns an existing topic or creates a new one
func (r *TopicRegistry) getOrCreateTopic(name string) *Topic {
	r.mu.RLock()
	t, ok := r.topics[name]
	r.mu.RUnlock()

	if ok {
		return t
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Double check
	t, ok = r.topics[name]
	if ok {
		return t
	}

	t = &Topic{
		Name:        name,
		Subscribers: make(map[string]bool),
		History:     NewHistoryBuffer(100), // Default capacity 100
	}
	r.topics[name] = t
	return t
}

// Subscribe adds a connection ID to a topic
func (r *TopicRegistry) Subscribe(topicName, connID string) {
	t := r.getOrCreateTopic(topicName)
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Subscribers[connID] = true
}

// Unsubscribe removes a connection ID from a topic
func (r *TopicRegistry) Unsubscribe(topicName, connID string) {
	r.mu.RLock()
	t, ok := r.topics[topicName]
	r.mu.RUnlock()

	if !ok {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.Subscribers, connID)

	// We could clean up empty topics here, but it might flap if clients
	// disconnect and reconnect frequently. A periodic cleanup is better.
}

// GetSubscribers returns a snapshot of connection IDs subscribed to a topic
func (r *TopicRegistry) GetSubscribers(topicName string) []string {
	r.mu.RLock()
	t, ok := r.topics[topicName]
	r.mu.RUnlock()

	if !ok {
		return nil
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	subs := make([]string, 0, len(t.Subscribers))
	for id := range t.Subscribers {
		subs = append(subs, id)
	}
	return subs
}
