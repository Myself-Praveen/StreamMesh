package grpc

import (
	"context"
	"testing"
	"time"
	"github.com/Myself-Praveen/StreamMesh/internal/grpc/pb"
	"github.com/Myself-Praveen/StreamMesh/internal/logger"
	"github.com/Myself-Praveen/StreamMesh/internal/redis"
)

func TestGRPCServerAndClient(t *testing.T) {
	logger.InitLogger("test")
	redis.CurrentNode.ID = "test-node-1"
	redis.CurrentNode.Hostname = "localhost"

	// Start Server
	server := NewServer()
	err := server.Start(50052)
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer server.Stop()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	// Test Client
	client, err := GlobalClientManager.GetClient("localhost:50052")
	if err != nil {
		t.Fatalf("Failed to get client: %v", err)
	}
	defer GlobalClientManager.CloseAll()

	// Test Gossip
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.GossipMessage{
		Sender: &pb.NodeInfo{
			Id:       "test-node-2",
			Hostname: "localhost",
		},
	}

	reply, err := client.Gossip(ctx, req)
	if err != nil {
		t.Fatalf("Gossip failed: %v", err)
	}

	if !reply.Success {
		t.Errorf("Expected success, got false")
	}

	// Test PublishMessage (assuming global dedup cache is nil in tests unless initialized)
	redis.GlobalDedupCache = &redis.DedupCache{}
	pubReq := &pb.PublishRequest{
		Topic:        "test-topic",
		Payload:      []byte("hello"),
		SenderNodeId: "test-node-2",
		EnvelopeId:   "env-1",
	}

	pubReply, err := client.PublishMessage(ctx, pubReq)
	if err != nil {
		t.Fatalf("PublishMessage failed: %v", err)
	}

	if !pubReply.Success {
		t.Errorf("Expected publish success, got false")
	}
}
