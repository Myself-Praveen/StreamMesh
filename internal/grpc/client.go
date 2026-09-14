package grpc

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Myself-Praveen/StreamMesh/internal/grpc/pb"
	"github.com/Myself-Praveen/StreamMesh/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ClientManager handles connections to other gRPC nodes
type ClientManager struct {
	conns map[string]*grpc.ClientConn
	mu    sync.RWMutex
}

var GlobalClientManager = &ClientManager{
	conns: make(map[string]*grpc.ClientConn),
}

// GetConnection returns an existing connection or creates a new one
func (m *ClientManager) GetConnection(address string) (*grpc.ClientConn, error) {
	m.mu.RLock()
	conn, exists := m.conns[address]
	m.mu.RUnlock()

	if exists {
		return conn, nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Double check
	if conn, exists = m.conns[address]; exists {
		return conn, nil
	}

	// Connect
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	
	if err != nil {
		logger.Log.Error("Failed to connect to gRPC node", zap.String("address", address), zap.Error(err))
		return nil, fmt.Errorf("failed to connect to node: %w", err)
	}

	m.conns[address] = conn
	return conn, nil
}

// CloseAll closes all active connections
func (m *ClientManager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for addr, conn := range m.conns {
		if err := conn.Close(); err != nil {
			logger.Log.Error("Failed to close gRPC connection", zap.String("address", addr), zap.Error(err))
		}
	}
	// clear map
	m.conns = make(map[string]*grpc.ClientConn)
}

// GetClient returns a ClusterService client for the given address
func (m *ClientManager) GetClient(address string) (pb.ClusterServiceClient, error) {
	conn, err := m.GetConnection(address)
	if err != nil {
		return nil, err
	}
	return pb.NewClusterServiceClient(conn), nil
}
