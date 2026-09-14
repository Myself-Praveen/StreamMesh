package pubsub

import (
	"testing"
)

func TestTopicRegistry(t *testing.T) {
	registry := NewTopicRegistry()
	
	registry.Subscribe("topic1", "client1")
	registry.Subscribe("topic1", "client2")
	registry.Subscribe("topic2", "client1")

	subs1 := registry.GetSubscribers("topic1")
	if len(subs1) != 2 {
		t.Errorf("Expected 2 subscribers for topic1, got %d", len(subs1))
	}

	subs2 := registry.GetSubscribers("topic2")
	if len(subs2) != 1 || subs2[0] != "client1" {
		t.Errorf("Expected 1 subscriber (client1) for topic2")
	}

	registry.Unsubscribe("topic1", "client1")
	subs1 = registry.GetSubscribers("topic1")
	if len(subs1) != 1 || subs1[0] != "client2" {
		t.Errorf("Expected 1 subscriber (client2) for topic1 after unsubscribe")
	}
}
