package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Marwan051/final_project_backend/internal/api/v1/handlers"
	"github.com/Marwan051/final_project_backend/internal/service/route_service"
)

// NewRouter returns a new router with all v1 API routes
func NewRouter(routingService route_service.Router) *http.ServeMux {
	mux := http.NewServeMux()

	// Create handler with the injected service
	routingHandler := handlers.NewRoutingHandler(routingService)

	// Health check
	mux.HandleFunc("GET /health", HealthHandler)

	// Routing endpoint
	mux.HandleFunc("POST /route", routingHandler.FindRoute)

	return mux
}

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// HealthHandler returns the health status of the service
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
