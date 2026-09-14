package grpc

import (
	"context"
	"time"

	"github.com/Myself-Praveen/StreamMesh/internal/grpc/pb"
	"github.com/Myself-Praveen/StreamMesh/internal/logger"
	"github.com/Myself-Praveen/StreamMesh/internal/redis"
	"go.uber.org/zap"
)

// Publisher handles broadcasting messages over gRPC
type Publisher struct {
	knownNodes map[string]string
}

var GlobalPublisher = &Publisher{
	knownNodes: make(map[string]string),
}

// BroadcastMessage sends a message to all known nodes via gRPC
func (p *Publisher) BroadcastMessage(topic string, payload []byte, envelopeID string) {
	req := &pb.PublishRequest{
		Topic:        topic,
		Payload:      payload,
		SenderNodeId: redis.CurrentNode.ID,
		EnvelopeId:   envelopeID,
	}

	for _, address := range p.knownNodes {
		go func(addr string) {
			client, err := GlobalClientManager.GetClient(addr)
			if err != nil {
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			reply, err := client.PublishMessage(ctx, req)
			if err != nil {
				logger.Log.Debug("gRPC Publish failed", zap.String("address", addr), zap.Error(err))
				return
			}
			
			if !reply.Success {
				logger.Log.Debug("gRPC Publish returned false success", zap.String("error", reply.Error))
			}
		}(address)
	}
}

// AddNode adds a target node for publishing
func (p *Publisher) AddNode(id, address string) {
	p.knownNodes[id] = address
}
