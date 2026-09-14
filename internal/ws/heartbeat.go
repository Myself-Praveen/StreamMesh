package ws

import (
	"time"

	"github.com/Myself-Praveen/StreamMesh/internal/config"
)

// HeartbeatConfig holds the configuration for WebSocket heartbeats
type HeartbeatConfig struct {
	WriteWait      time.Duration
	PongWait       time.Duration
	PingPeriod     time.Duration
	MaxMessageSize int64
}

// GetHeartbeatConfig generates the heartbeat config from global app config
func GetHeartbeatConfig() HeartbeatConfig {
	interval := config.AppConfig.WebSocket.HeartbeatInterval
	if interval <= 0 {
		interval = 30 // default 30s
	}

	maxSize := config.AppConfig.WebSocket.MaxMessageSize
	if maxSize <= 0 {
		maxSize = 4096 // default 4KB
	}

	pongWait := time.Duration(interval) * time.Second
	pingPeriod := (pongWait * 9) / 10

	return HeartbeatConfig{
		WriteWait:      10 * time.Second,
		PongWait:       pongWait,
		PingPeriod:     pingPeriod,
		MaxMessageSize: int64(maxSize),
	}
}
