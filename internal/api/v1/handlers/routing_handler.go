package handlers

import (
	"log"
	"net/http"

	"github.com/Marwan051/final_project_backend/internal/service/route_service"
	"github.com/Marwan051/final_project_backend/internal/utils"
)

type RoutingHandler struct {
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

	if err := utils.DecodeJSONBody(r, &req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Validate required fields
	if req.StartLat == 0 || req.StartLon == 0 || req.EndLat == 0 || req.EndLon == 0 {
		utils.WriteJSONError(w, http.StatusBadRequest, "Missing required coordinates")
		return
	}

	// Call the routing service
	resp, err := h.routerService.FindRoute(r.Context(), req)
	if err != nil {
		log.Printf("Error finding route: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to find route")
		return
	}

	// Return JSON response
	if err := utils.WriteJSONResponse(w, http.StatusOK, resp); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
