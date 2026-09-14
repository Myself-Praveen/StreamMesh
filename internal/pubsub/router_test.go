package pubsub

import (
	"testing"

	"github.com/Myself-Praveen/StreamMesh/internal/ws"
	"github.com/Myself-Praveen/StreamMesh/internal/logger"
)

func init() {
	logger.InitLogger("development")
}

func TestRouter_Publish(t *testing.T) {
	manager := ws.NewManager()
	registry := NewTopicRegistry()
	router := NewRouter(registry, manager)

	conn1 := ws.NewConnection("c1", nil)
	manager.Add(conn1)
	registry.Subscribe("test.topic", "c1")

	err := router.Publish("test.topic", []byte(`{"data":"test"}`))
	if err != nil {
		t.Fatalf("Failed to publish: %v", err)
	}

	// Read from conn1 Send channel
	select {
	case msg := <-conn1.Send:
		if len(msg) == 0 {
			t.Errorf("Expected message, got empty bytes")
		}
	default:
		t.Errorf("Expected message in channel, got none")
	}
}
