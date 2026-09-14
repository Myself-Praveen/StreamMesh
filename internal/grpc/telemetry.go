package grpc

import (
	"context"
	"time"

	"github.com/Myself-Praveen/StreamMesh/internal/grpc/pb"
	"github.com/Myself-Praveen/StreamMesh/internal/logger"
	"github.com/Myself-Praveen/StreamMesh/internal/redis"
	"go.uber.org/zap"
)

// TelemetryManager handles relaying telemetry to other nodes
type TelemetryManager struct {
	knownNodes map[string]string
}

var GlobalTelemetryManager = &TelemetryManager{
	knownNodes: make(map[string]string),
}

// StartTelemetryLoop starts sending telemetry periodically
func (m *TelemetryManager) StartTelemetryLoop(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				m.broadcastTelemetry(ctx)
			}
		}
	}()
}

// broadcastTelemetry sends telemetry to all known nodes
func (m *TelemetryManager) broadcastTelemetry(ctx context.Context) {
	for _, address := range m.knownNodes {
		go func(addr string) {
			client, err := GlobalClientManager.GetClient(addr)
			if err != nil {
				return
			}
			
			streamCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			
			stream, err := client.StreamTelemetry(streamCtx)
			if err != nil {
				logger.Log.Debug("Failed to open telemetry stream", zap.Error(err))
				return
			}
			
			// Get actual stats in real app, these are placeholders
			activeConns := int64(100)
			msgSec := 50.5
			p99Lat := 1.2
			
			data := &pb.TelemetryData{
				NodeId:            redis.CurrentNode.ID,
				ActiveConnections: activeConns,
				MessagesPerSecond: msgSec,
				P99LatencyMs:      p99Lat,
			}
			
			if err := stream.Send(data); err != nil {
				logger.Log.Debug("Failed to send telemetry", zap.Error(err))
				return
			}
			
			if _, err := stream.CloseAndRecv(); err != nil {
				logger.Log.Debug("Error closing telemetry stream", zap.Error(err))
			}
		}(address)
	}
}

// AddNode adds a known node to telemetry stream targets
func (m *TelemetryManager) AddNode(id, address string) {
	m.knownNodes[id] = address
}
