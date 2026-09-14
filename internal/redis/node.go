package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Myself-Praveen/StreamMesh/internal/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Node represents a server instance in the cluster
type Node struct {
	ID        string `json:"id"`
	Hostname  string `json:"hostname"`
	StartTime time.Time `json:"start_time"`
}

var CurrentNode Node

func init() {
	CurrentNode = Node{
		ID:        uuid.New().String(),
		Hostname:  "streammesh-node", // could be populated from os.Hostname
		StartTime: time.Now(),
	}
}

const NodeRegistryKey = "streammesh:nodes"

// RegisterNode registers the current node in Redis with a TTL
func RegisterNode(ctx context.Context) {
	if Client == nil {
		return
	}

	data, _ := json.Marshal(CurrentNode)

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			// Register with 30s TTL
			err := Client.HSet(ctx, NodeRegistryKey, CurrentNode.ID, data).Err()
			if err != nil {
				logger.Log.Error("Failed to register node", zap.Error(err))
			} else {
				// Refresh TTL for the hash field would require a separate string key per node usually,
				// but for simplicity, we just keep writing to the hash. 
				// To prune dead nodes, we need a separate mechanism or use individual keys with TTL.
				// Let's use individual keys for TTL support
				key := "streammesh:node:" + CurrentNode.ID
				Client.Set(ctx, key, data, 30*time.Second)
			}

			select {
			case <-ctx.Done():
				// Cleanup on exit
				key := "streammesh:node:" + CurrentNode.ID
				Client.Del(context.Background(), key)
				return
			case <-ticker.C:
			}
		}
	}()
}
