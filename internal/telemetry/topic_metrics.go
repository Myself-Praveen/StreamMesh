package telemetry

import (
	"sync"
	"sync/atomic"
)

// TopicMetric tracks metrics for a specific topic
type TopicMetric struct {
	MessagesPublished int64
	MessagesDelivered int64
	Subscribers       int32
}

// TopicMetricsTracker tracks metrics across all topics
type TopicMetricsTracker struct {
	mu     sync.RWMutex
	topics map[string]*TopicMetric
}

// GlobalTopicTracker is the singleton for topic metrics
var GlobalTopicTracker = &TopicMetricsTracker{
	topics: make(map[string]*TopicMetric),
}

// getOrCreateMetric returns the metric for a topic, creating it if needed
func (t *TopicMetricsTracker) getOrCreateMetric(topic string) *TopicMetric {
	t.mu.RLock()
	m, ok := t.topics[topic]
	t.mu.RUnlock()
	
	if ok {
		return m
	}
	
	t.mu.Lock()
	defer t.mu.Unlock()
	
	if m, ok = t.topics[topic]; ok {
		return m
	}
	
	m = &TopicMetric{}
	t.topics[topic] = m
	return m
}

// RecordPublish increments published count for a topic
func (t *TopicMetricsTracker) RecordPublish(topic string) {
	m := t.getOrCreateMetric(topic)
	atomic.AddInt64(&m.MessagesPublished, 1)
}

// RecordDelivery increments delivered count for a topic
func (t *TopicMetricsTracker) RecordDelivery(topic string, count int) {
	m := t.getOrCreateMetric(topic)
	atomic.AddInt64(&m.MessagesDelivered, int64(count))
}

// SetSubscribers sets the subscriber count for a topic
func (t *TopicMetricsTracker) SetSubscribers(topic string, count int) {
	m := t.getOrCreateMetric(topic)
	atomic.StoreInt32(&m.Subscribers, int32(count))
}

// Snapshot returns the current metrics for all topics
func (t *TopicMetricsTracker) Snapshot() map[string]map[string]interface{} {
	t.mu.RLock()
	defer t.mu.RUnlock()
	
	result := make(map[string]map[string]interface{})
	for name, m := range t.topics {
		result[name] = map[string]interface{}{
			"messages_published": atomic.LoadInt64(&m.MessagesPublished),
			"messages_delivered": atomic.LoadInt64(&m.MessagesDelivered),
			"subscribers":        atomic.LoadInt32(&m.Subscribers),
		}
	}
	return result
}
