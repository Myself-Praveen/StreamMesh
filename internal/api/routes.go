package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Myself-Praveen/StreamMesh/internal/logger"
	"github.com/Myself-Praveen/StreamMesh/internal/ws"
	"go.uber.org/zap"
)

// SetupRoutes configures the basic HTTP and WebSocket routes
func SetupRoutes(manager *ws.Manager) *http.ServeMux {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// WebSocket upgrade endpoint
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := ws.UpgradeHandler(w, r)
		if err != nil {
			logger.Log.Error("Failed to upgrade connection", zap.Error(err))
			return
		}

		// Generate a simple ID for now, later use UUIDs
		connID := fmt.Sprintf("conn-%d", time.Now().UnixNano())
		client := ws.NewConnection(connID, conn)

		manager.Add(client)
		logger.Log.Info("Client connected", zap.String("id", connID))

		// Start pump goroutines
		go client.WritePump()
		go client.ReadPump(
			func(c *ws.Connection) {
				manager.Remove(c.ID)
				logger.Log.Info("Client disconnected", zap.String("id", c.ID))
			},
			func(c *ws.Connection, msg []byte) {
				// Process incoming message
				// For now, just log it. Later, route to pub/sub engine.
				logger.Log.Debug("Received message", zap.String("id", c.ID), zap.ByteString("msg", msg))
			},
		)
	})

	return mux
}
