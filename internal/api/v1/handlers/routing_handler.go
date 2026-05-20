package handlers

import (
	"log"
	"net/http"

	route_service "github.com/Marwan051/final_project_backend/internal/service/routing"
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

// FindRoute finds the best route between two points
// @Summary      Find Route
// @Description  Calculates the best route from start coordinates to end coordinates
// @Tags         routing
// @Accept       json
// @Produce      json
// @Param        request body route_service.RouteRequest true "Route Request"
// @Success      200  {object}  route_service.RouteResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/v1/route [post]
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

// ReloadPrefixTimes reloads prefix times for routing data.
// @Summary      Reload Prefix Times
// @Description  Admin only. Reloads prefix times for the routing graph.
// @Tags         routing
// @Accept       json
// @Produce      json
// @Success      200  {object}  route_service.AdminOperationResponse
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/v1/routing/reload-prefix-times [post]
func (h *RoutingHandler) ReloadPrefixTimes(w http.ResponseWriter, r *http.Request) {
	resp, err := h.routerService.ReloadPrefixTimes(r.Context())
	if err != nil {
		log.Printf("Error reloading prefix times: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to reload prefix times")
		return
	}

	if err := utils.WriteJSONResponse(w, http.StatusOK, resp); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// RebuildNetwork rebuilds the routing network.
// @Summary      Rebuild Routing Network
// @Description  Admin only. Rebuilds the routing network from source data.
// @Tags         routing
// @Accept       json
// @Produce      json
// @Success      200  {object}  route_service.AdminOperationResponse
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/v1/routing/rebuild-network [post]
func (h *RoutingHandler) RebuildNetwork(w http.ResponseWriter, r *http.Request) {
	resp, err := h.routerService.RebuildNetwork(r.Context())
	if err != nil {
		log.Printf("Error rebuilding network: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to rebuild network")
		return
	}

	if err := utils.WriteJSONResponse(w, http.StatusOK, resp); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
