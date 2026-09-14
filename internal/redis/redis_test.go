package redis

import (
	"context"
	"testing"
)

func TestDedupCache(t *testing.T) {
	cache := &DedupCache{}

	// First time should be false
	if cache.CheckAndAdd("msg1") {
		t.Errorf("Expected msg1 to be new")
	}

	// Second time should be true (duplicate)
	if !cache.CheckAndAdd("msg1") {
		t.Errorf("Expected msg1 to be duplicate")
	}

	// Different ID should be new
	if cache.CheckAndAdd("msg2") {
		t.Errorf("Expected msg2 to be new")
	}
}

func TestRedisIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Very basic check - we can't test properly without a real config/Redis
	// In a real environment, we would use a library like miniredis or testcontainers
	// For now, we just skip.
	t.Skip("Requires running Redis instance")
}

func TestNodeRegistration(t *testing.T) {
	if CurrentNode.ID == "" {
		t.Errorf("CurrentNode ID should be generated on init")
	}

	if CurrentNode.Hostname == "" {
		t.Errorf("CurrentNode Hostname should be populated")
	}
}

func TestTopicSyncMessage(t *testing.T) {
	// Just testing the struct and JSON marshaling implicitly via BroadcastTopicSync
	// Since Client is nil in tests, it should NOOP
	BroadcastTopicSync(context.Background(), "test-topic", "subscribe")
}
