package grpc

import (
	"context"
	"time"

	"github.com/Myself-Praveen/StreamMesh/internal/grpc/pb"
	"github.com/Myself-Praveen/StreamMesh/internal/logger"
	"github.com/Myself-Praveen/StreamMesh/internal/redis"
	"go.uber.org/zap"
)

// GossipManager handles background gossip protocol
type GossipManager struct {
	server       *Server
	knownNodes   map[string]string // ID -> Address
}

var GlobalGossipManager = &GossipManager{
	knownNodes: make(map[string]string),
}

func init() {
	// Initialize with some dummy or seed nodes in a real scenario
}

// StartGossipLoop starts periodically gossiping with known nodes
func (m *GossipManager) StartGossipLoop(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				m.GossipWithNodes()
			}
		}
	}()
}

// GossipWithNodes picks random nodes and gossips
func (m *GossipManager) GossipWithNodes() {
	// In a real implementation, we'd fetch nodes from Redis or local list,
	// pick a few random ones and send a gossip message.
	for _, address := range m.knownNodes {
		client, err := GlobalClientManager.GetClient(address)
		if err != nil {
			logger.Log.Debug("Failed to get client for gossip", zap.String("address", address), zap.Error(err))
			continue
		}

		req := &pb.GossipMessage{
			Sender: &pb.NodeInfo{
				Id:       redis.CurrentNode.ID,
				Hostname: redis.CurrentNode.Hostname,
			},
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		reply, err := client.Gossip(ctx, req)
		cancel()

		if err != nil {
			logger.Log.Debug("Gossip failed", zap.String("address", address), zap.Error(err))
			continue
		}

		if reply.Success {
			logger.Log.Debug("Gossip successful", zap.String("address", address))
		}
	}
}

// AddNode adds a known node to gossip with
func (m *GossipManager) AddNode(id, address string) {
	m.knownNodes[id] = address
}
