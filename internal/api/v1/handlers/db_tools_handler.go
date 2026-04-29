package handlers

import (
	"log"
	"net/http"

	"github.com/Marwan051/final_project_backend/internal/service/db_tools"
	pb "github.com/Marwan051/final_project_backend/internal/service/db_tools/proto"
	"github.com/Marwan051/final_project_backend/internal/utils"
)

type DbToolsHandler struct {
	client db_tools.DbTools
}

func NewDbToolsHandler(client db_tools.DbTools) *DbToolsHandler {
	return &DbToolsHandler{
		client: client,
	}
}

// NearbyTrips gets nearby trips
// @Summary      Get Nearby Trips
// @Description  Finds trips near a given location within a specified radius
// @Tags         db_tools
// @Accept       json
// @Produce      json
// @Param        request body pb.NearbyTripsRequest true "Nearby Trips Request"
// @Success      200  {object}  pb.NearbyTripsResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/v1/nearby-trips [post]
func (h *DbToolsHandler) NearbyTrips(w http.ResponseWriter, r *http.Request) {
	var req pb.NearbyTripsRequest
	if err := utils.DecodeJSONBody(r, &req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	resp, err := h.client.NearbyTrips(r.Context(), &req)
	if err != nil {
		log.Printf("Error getting nearby trips: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to get nearby trips: "+err.Error())
		return
	}

	if err := utils.WriteJSONResponse(w, http.StatusOK, resp); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
