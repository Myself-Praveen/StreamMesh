package ws

import (
	"context"
	"time"

	"github.com/Myself-Praveen/StreamMesh/internal/logger"
	"go.uber.org/zap"
)

// Pruner periodically scans for zombie connections and removes them
type Pruner struct {
	manager  *Manager
	interval time.Duration
	timeout  time.Duration
}

// NewPruner creates a new zombie connection pruner
func NewPruner(manager *Manager, interval, timeout time.Duration) *Pruner {
	return &Pruner{
		manager:  manager,
		interval: interval,
		timeout:  timeout,
	}
}

// Start runs the pruning loop until the context is canceled
func (p *Pruner) Start(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Stopping zombie connection pruner")
			return
		case <-ticker.C:
			p.prune()
		}
	}
}

func (p *Pruner) prune() {
	now := time.Now()
	connections := p.manager.GetAll()

	var pruned int
	for _, conn := range connections {
		if now.Sub(conn.GetLastPing()) > p.timeout {
			logger.Log.Warn("Pruning zombie connection", zap.String("id", conn.ID))
			if conn.Conn != nil {
				conn.Conn.Close() // This will trigger the read pump to return, unregistering the conn
			}
			pruned++
		}
	}

	if pruned > 0 {
		logger.Log.Info("Pruned zombie connections", zap.Int("count", pruned))
	}
}
