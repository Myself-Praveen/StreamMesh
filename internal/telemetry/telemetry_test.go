package telemetry

import (
	"testing"
)

func TestCollector(t *testing.T) {
	GlobalCollector.IncActiveConnections()
	GlobalCollector.RecordPublish()
	GlobalCollector.RecordDelivery()
	GlobalCollector.RecordBytesSent(100)
	GlobalCollector.RecordBytesReceived(200)

	snapshot := GlobalCollector.Snapshot()

	if snapshot["active_connections"] != 1 {
		t.Errorf("Expected 1 active connection, got %d", snapshot["active_connections"])
	}
	if snapshot["messages_published"] != 1 {
		t.Errorf("Expected 1 published, got %d", snapshot["messages_published"])
	}
	if snapshot["messages_delivered"] != 1 {
		t.Errorf("Expected 1 delivered, got %d", snapshot["messages_delivered"])
	}
	if snapshot["bytes_sent"] != 100 {
		t.Errorf("Expected 100 bytes sent, got %d", snapshot["bytes_sent"])
	}
	if snapshot["bytes_received"] != 200 {
		t.Errorf("Expected 200 bytes received, got %d", snapshot["bytes_received"])
	}

	GlobalCollector.DecActiveConnections()
	snapshot = GlobalCollector.Snapshot()
	if snapshot["active_connections"] != 0 {
		t.Errorf("Expected 0 active connections after decrement, got %d", snapshot["active_connections"])
	}
}

func TestHistogram(t *testing.T) {
	h := NewHistogram(10)

	// Record values 1 through 10
	for i := 1; i <= 10; i++ {
		h.Record(float64(i))
	}

	// Calculate percentiles
	p50 := h.Percentile(50) // should be 6
	p90 := h.Percentile(90) // should be 10
	p99 := h.Percentile(99) // should be 10 (since 10 is max)

	if p50 != 6 {
		t.Errorf("Expected p50 to be 6, got %v", p50)
	}
	if p90 != 10 {
		t.Errorf("Expected p90 to be 10, got %v", p90)
	}
	if p99 != 10 {
		t.Errorf("Expected p99 to be 10, got %v", p99)
	}
}

func TestTopicMetrics(t *testing.T) {
	GlobalTopicTracker.RecordPublish("test-topic")
	GlobalTopicTracker.RecordDelivery("test-topic", 5)
	GlobalTopicTracker.SetSubscribers("test-topic", 3)

	snapshot := GlobalTopicTracker.Snapshot()
	stats, ok := snapshot["test-topic"]
	if !ok {
		t.Fatalf("test-topic metrics not found")
	}

	if stats["messages_published"].(int64) != 1 {
		t.Errorf("Expected 1 published, got %v", stats["messages_published"])
	}
	if stats["messages_delivered"].(int64) != 5 {
		t.Errorf("Expected 5 delivered, got %v", stats["messages_delivered"])
	}
	if stats["subscribers"].(int32) != 3 {
		t.Errorf("Expected 3 subscribers, got %v", stats["subscribers"])
	}
}
