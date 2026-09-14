package api

import (
	"encoding/json"
	"net/http"
	
	"github.com/Myself-Praveen/StreamMesh/internal/pubsub"
	"github.com/Myself-Praveen/StreamMesh/internal/redis"
	"github.com/Myself-Praveen/StreamMesh/internal/ws"
)

// AdminHandler encapsulates the dependencies for admin endpoints
type AdminHandler struct {
	manager  *ws.Manager
	registry *pubsub.TopicRegistry
	router   *pubsub.Router
}

// NewAdminHandler creates a new AdminHandler
func NewAdminHandler(manager *ws.Manager, registry *pubsub.TopicRegistry, router *pubsub.Router) *AdminHandler {
	return &AdminHandler{
		manager:  manager,
		registry: registry,
		router:   router,
	}
}

// HandleChannels list channels and subscribers
func (h *AdminHandler) HandleChannels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	topics := h.registry.GetAllTopics()
	
	response := map[string]interface{}{
		"channels": topics,
		"count":    len(topics),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleConnections lists active connections
func (h *AdminHandler) HandleConnections(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	conns := h.manager.GetAll()
	
	response := make([]map[string]interface{}, 0, len(conns))
	for _, c := range conns {
		userID, _ := c.GetMetadata("user_id")
		ip, _ := c.GetMetadata("ip")
		
		response = append(response, map[string]interface{}{
			"id":         c.ID,
			"user_id":    userID,
			"ip_address": ip,
			"connected":  c.CreatedAt,
			"metadata":   c.Metadata,
		})
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"connections": response,
		"count":       len(response),
	})
}

// HandleCluster lists cluster nodes
func (h *AdminHandler) HandleCluster(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	nodes := redis.GetActiveNodes(r.Context())
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"nodes": nodes,
		"count": len(nodes),
		"self":  redis.CurrentNode.ID,
	})
}

// HandlePublish publishes a message via REST
func (h *AdminHandler) HandlePublish(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		Topic   string          `json:"topic"`
		Payload json.RawMessage `json:"payload"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if req.Topic == "" {
		http.Error(w, "Topic is required", http.StatusBadRequest)
		return
	}
	
	h.router.Publish(req.Topic, req.Payload)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "published", "topic": req.Topic})
}
