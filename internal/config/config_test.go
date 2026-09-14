package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a dummy config file for testing in the current directory
	os.WriteFile("default.yaml", []byte(`
server:
  port: "9999"
  env: "test"
redis:
  url: "redis:6379"
websocket:
  heartbeat_interval: 10
  max_message_size: 1024
`), 0644)
	defer os.Remove("default.yaml")

	// Set env var to override port
	os.Setenv("SERVER_PORT", "8888")
	defer os.Unsetenv("SERVER_PORT")

	err := LoadConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if AppConfig == nil {
		t.Fatalf("expected AppConfig to be initialized")
	}

	if AppConfig.Server.Port != "8888" {
		t.Errorf("expected port 8888, got %v", AppConfig.Server.Port)
	}
	if AppConfig.Server.Env != "test" {
		t.Errorf("expected env test, got %v", AppConfig.Server.Env)
	}
	if AppConfig.WebSocket.HeartbeatInterval != 10 {
		t.Errorf("expected heartbeat 10, got %v", AppConfig.WebSocket.HeartbeatInterval)
	}
}
