package telemetry

import (
	"sync/atomic"
)

// Collector tracks global application metrics
type Collector struct {
	ActiveConnections int64
	MessagesPublished int64
	MessagesDelivered int64
	BytesSent         int64
	BytesReceived     int64
	Errors            int64
}

// GlobalCollector is the singleton collector for the node
var GlobalCollector = &Collector{}

// IncActiveConnections increments the active connection count
func (c *Collector) IncActiveConnections() {
	atomic.AddInt64(&c.ActiveConnections, 1)
}

// DecActiveConnections decrements the active connection count
func (c *Collector) DecActiveConnections() {
	atomic.AddInt64(&c.ActiveConnections, -1)
}

// RecordPublish records a published message
func (c *Collector) RecordPublish() {
	atomic.AddInt64(&c.MessagesPublished, 1)
}

// RecordDelivery records a delivered message
func (c *Collector) RecordDelivery() {
	atomic.AddInt64(&c.MessagesDelivered, 1)
}

// RecordBytesSent records bytes sent
func (c *Collector) RecordBytesSent(bytes int64) {
	atomic.AddInt64(&c.BytesSent, bytes)
}

// RecordBytesReceived records bytes received
func (c *Collector) RecordBytesReceived(bytes int64) {
	atomic.AddInt64(&c.BytesReceived, bytes)
}

// RecordError records an error
func (c *Collector) RecordError() {
	atomic.AddInt64(&c.Errors, 1)
}

// Snapshot returns a copy of current metrics
func (c *Collector) Snapshot() map[string]int64 {
	return map[string]int64{
		"active_connections": atomic.LoadInt64(&c.ActiveConnections),
		"messages_published": atomic.LoadInt64(&c.MessagesPublished),
		"messages_delivered": atomic.LoadInt64(&c.MessagesDelivered),
		"bytes_sent":         atomic.LoadInt64(&c.BytesSent),
		"bytes_received":     atomic.LoadInt64(&c.BytesReceived),
		"errors":             atomic.LoadInt64(&c.Errors),
	}
}
