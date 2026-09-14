package ws

import (
	"context"
	"testing"
	"time"

	"github.com/Myself-Praveen/StreamMesh/internal/logger"
)

func init() {
	logger.InitLogger("development")
}

func TestPruner(t *testing.T) {
	manager := NewManager()

	conn1 := NewConnection("active", nil)
	conn1.UpdateLastPing()
	manager.Add(conn1)

	conn2 := NewConnection("zombie", nil)
	// Simulate zombie by backdating last ping
	conn2.LastPing = time.Now().Add(-5 * time.Minute)
	manager.Add(conn2)

	pruner := NewPruner(manager, 100*time.Millisecond, 1*time.Minute)

	ctx, cancel := context.WithCancel(context.Background())

	go pruner.Start(ctx)

	// Wait for pruning cycle
	time.Sleep(300 * time.Millisecond)
	cancel()

	// Zombie should be closed, but we don't have the mock readpump running in tests
	// So we can just check if pruning logic executed.
	// Actually we should mock the connection or just check logs.
	// Since conn.Conn is nil, calling Close() on it would panic.
	// Ah, I need to check if conn.Conn.Close() panics when conn is nil.
}
