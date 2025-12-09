package handlers

import (
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
	// ... decode JSON ...

	// Call the interface method
	journeys, err := h.routerService.FindRoute(r.Context(), params)

	// ...
}
