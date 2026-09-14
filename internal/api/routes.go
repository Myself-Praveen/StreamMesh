package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Myself-Praveen/StreamMesh/internal/auth"
	"github.com/Myself-Praveen/StreamMesh/internal/config"
	"github.com/Myself-Praveen/StreamMesh/internal/logger"
	"github.com/Myself-Praveen/StreamMesh/internal/ratelimit"
	"github.com/Myself-Praveen/StreamMesh/internal/ws"
	"go.uber.org/zap"
)

// SetupRoutes configures the basic HTTP and WebSocket routes
func SetupRoutes(manager *ws.Manager, msgHandler *ws.MessageHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
	
	// Metrics endpoints
	mux.HandleFunc("/api/metrics", MetricsHandler)
	mux.HandleFunc("/api/metrics/stream", MetricsStreamHandler)

	ipLimiter := ratelimit.NewIPRateLimiter(10, 2.0) // Allow 10 burst, 2 per sec

	// WebSocket upgrade endpoint
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ip := ratelimit.GetIP(r)
		if !ipLimiter.Allow(ip) {
			logger.Log.Warn("IP connection rate limit exceeded", zap.String("ip", ip))
			http.Error(w, "Too many connection attempts", http.StatusTooManyRequests)
			return
		}

		if config.AppConfig.Auth.Enabled {
			tokenString := r.URL.Query().Get("token")
			apiKeyString := r.URL.Query().Get("api_key")
			
			authorized := false
			
			// Check JWT token
			if tokenString != "" {
				validator := auth.NewJWTValidator(config.AppConfig.Auth.JWTSecret)
				claims, err := validator.ValidateToken(tokenString)
				if err == nil {
					authorized = true
					_ = claims // Could attach to context
				} else {
					logger.Log.Debug("JWT validation failed", zap.Error(err), zap.String("ip", ip))
				}
			}
			
			// Check API key if not authorized yet
			if !authorized && apiKeyString != "" {
				keyValidator := auth.NewAPIKeyValidator(config.AppConfig.Auth.APIKeys)
				if keyValidator.ValidateKey(apiKeyString) {
					authorized = true
				} else {
					logger.Log.Debug("API Key validation failed", zap.String("ip", ip))
				}
			}

			if !authorized {
				logger.Log.Warn("Unauthorized access attempt", zap.String("ip", ip))
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

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
				msgHandler.HandleMessage(c, msg)
			},
		)
	})

	return mux
}
