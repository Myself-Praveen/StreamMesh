package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/Myself-Praveen/StreamMesh/internal/grpc/pb"
	"github.com/Myself-Praveen/StreamMesh/internal/logger"
	"github.com/Myself-Praveen/StreamMesh/internal/pubsub"
	"github.com/Myself-Praveen/StreamMesh/internal/redis"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server implements the gRPC ClusterService
type Server struct {
	pb.UnimplementedClusterServiceServer
	grpcServer *grpc.Server
}

// NewServer creates a new gRPC cluster server
func NewServer() *Server {
	return &Server{
		grpcServer: grpc.NewServer(),
	}
}

// Start starts the gRPC server on the specified port
func (s *Server) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	pb.RegisterClusterServiceServer(s.grpcServer, s)
	reflection.Register(s.grpcServer)

	logger.Log.Info("Starting gRPC cluster server", zap.Int("port", port))
	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			logger.Log.Error("gRPC server failed", zap.Error(err))
		}
	}()

	return nil
}

// Stop gracefully stops the gRPC server
func (s *Server) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}

// Gossip handles gossip messages from other nodes
func (s *Server) Gossip(ctx context.Context, req *pb.GossipMessage) (*pb.GossipReply, error) {
	// Accept gossip and potentially update local routing tables
	logger.Log.Debug("Received gossip", zap.String("from", req.Sender.Id))
	
	reply := &pb.GossipReply{
		Success: true,
		KnownNodes: []*pb.NodeInfo{
			{
				Id:       redis.CurrentNode.ID,
				Hostname: redis.CurrentNode.Hostname,
			},
		},
	}
	
	return reply, nil
}

// StreamTelemetry accepts telemetry data streams
func (s *Server) StreamTelemetry(stream pb.ClusterService_StreamTelemetryServer) error {
	for {
		data, err := stream.Recv()
		if err != nil {
			return err
		}
		
		logger.Log.Debug("Received telemetry", 
			zap.String("node_id", data.NodeId),
			zap.Float64("msgs_sec", data.MessagesPerSecond))
			
		err = stream.SendAndClose(&pb.TelemetryAck{Received: true})
		if err != nil {
			return err
		}
	}
}

// PublishMessage handles direct cross-node publishing
func (s *Server) PublishMessage(ctx context.Context, req *pb.PublishRequest) (*pb.PublishReply, error) {
	// Deduplicate if we saw this message already
	if redis.GlobalDedupCache.CheckAndAdd(req.EnvelopeId) {
		return &pb.PublishReply{Success: true}, nil
	}

	// Publish locally
	if pubsub.GlobalRouter != nil {
		pubsub.GlobalRouter.PublishLocal(req.Topic, req.Payload)
	}

	return &pb.PublishReply{Success: true}, nil
}
