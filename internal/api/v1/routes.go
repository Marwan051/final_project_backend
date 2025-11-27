package v1

import (
	"encoding/json"
	"net/http"
	"time"
)

// NewRouter returns a new router with all v1 API routes
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", HealthHandler)

	// Add more routes here
	// mux.HandleFunc("GET /routes", handlers.GetRoutes)
	// mux.HandleFunc("POST /routes", handlers.CreateRoute)

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
