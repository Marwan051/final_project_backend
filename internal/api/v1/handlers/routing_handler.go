package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Marwan051/final_project_backend/internal/service/route_service"
)

// internal/api/v1/handlers/routing_handler.go

type RoutingHandler struct {
	// DEPENDENCY INJECTION: We ask for the Interface
	routerService route_service.Router
}

// Constructor accepts the interface
func NewRoutingHandler(router route_service.Router) *RoutingHandler {
	return &RoutingHandler{
		routerService: router,
	}
}

func (h *RoutingHandler) FindRoute(w http.ResponseWriter, r *http.Request) {
	var req route_service.RouteRequest

	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.StartLat == 0 || req.StartLon == 0 || req.EndLat == 0 || req.EndLon == 0 {
		http.Error(w, "Missing required coordinates", http.StatusBadRequest)
		return
	}

	// Call the routing service
	resp, err := h.routerService.FindRoute(r.Context(), req)
	if err != nil {
		log.Printf("Error finding route: %v", err)
		http.Error(w, "Failed to find route: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func (h *RoutingHandler) GrpcHealthHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := h.routerService.HealthCheck(r.Context())
	if err != nil {
		log.Printf("could not complete grpc health check")
		http.Error(w, "Failed to complete grpc health check: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
