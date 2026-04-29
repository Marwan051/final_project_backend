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
