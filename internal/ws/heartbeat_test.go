package ws

import (
	"testing"
	"time"

	"github.com/Myself-Praveen/StreamMesh/internal/config"
)

func TestGetHeartbeatConfig(t *testing.T) {
	// Mock config
	config.AppConfig = &config.Config{}
	config.AppConfig.WebSocket.HeartbeatInterval = 20
	config.AppConfig.WebSocket.MaxMessageSize = 1024

	hb := GetHeartbeatConfig()
	if hb.PongWait != 20*time.Second {
		t.Errorf("Expected PongWait 20s, got %v", hb.PongWait)
	}
	if hb.MaxMessageSize != 1024 {
		t.Errorf("Expected MaxMessageSize 1024, got %v", hb.MaxMessageSize)
	}
}
