package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	
	"github.com/Myself-Praveen/StreamMesh/internal/telemetry"
)

// MetricsHandler handles HTTP requests for system metrics
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	globalStats := telemetry.GlobalCollector.Snapshot()
	topicStats := telemetry.GlobalTopicTracker.Snapshot()
	latencies := telemetry.GlobalLatencyHistogram.Snapshots()
	sysStats := telemetry.GetSystemMetrics()
	
	response := map[string]interface{}{
		"global":    globalStats,
		"topics":    topicStats,
		"latencies": latencies,
		"system":    sysStats,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// MetricsStreamHandler streams metrics over Server-Sent Events (SSE) or WebSocket
func MetricsStreamHandler(w http.ResponseWriter, r *http.Request) {
	// Simple SSE for streaming metrics
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}
	
	ctx := r.Context()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			globalStats := telemetry.GlobalCollector.Snapshot()
			topicStats := telemetry.GlobalTopicTracker.Snapshot()
			latencies := telemetry.GlobalLatencyHistogram.Snapshots()
			sysStats := telemetry.GetSystemMetrics()
			
			response := map[string]interface{}{
				"global":    globalStats,
				"topics":    topicStats,
				"latencies": latencies,
				"system":    sysStats,
			}
			
			b, _ := json.Marshal(response)
			fmt.Fprintf(w, "data: %s\n\n", b)
			flusher.Flush()
		}
	}
}
