package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Myself-Praveen/StreamMesh/internal/pubsub"
	"github.com/Myself-Praveen/StreamMesh/internal/ws"
)

func TestAdminHandlerChannels(t *testing.T) {
	manager := ws.NewManager()
	registry := pubsub.NewTopicRegistry()
	router := pubsub.NewRouter(registry, manager)
	
	registry.Subscribe("test-topic", "conn-1")
	
	handler := NewAdminHandler(manager, registry, router)
	
	req := httptest.NewRequest(http.MethodGet, "/api/admin/channels", nil)
	w := httptest.NewRecorder()
	
	handler.HandleChannels(w, req)
	
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Result().StatusCode)
	}
	
	var res map[string]interface{}
	json.NewDecoder(w.Body).Decode(&res)
	
	if res["count"].(float64) != 1 {
		t.Errorf("Expected 1 channel, got %v", res["count"])
	}
}
