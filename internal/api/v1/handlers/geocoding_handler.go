package handlers

import (
	"log"
	"net/http"

	"github.com/Marwan051/final_project_backend/internal/service/geocoding"
	pb "github.com/Marwan051/final_project_backend/internal/service/geocoding/proto"
	"github.com/Marwan051/final_project_backend/internal/utils"
)

type GeocodingHandler struct {
	client geocoding.Geocoding
}

func NewGeocodingHandler(client geocoding.Geocoding) *GeocodingHandler {
	return &GeocodingHandler{
		client: client,
	}
}

func (h *GeocodingHandler) Geocode(w http.ResponseWriter, r *http.Request) {
	var req pb.GeocodeRequest
	if err := utils.DecodeJSONBody(r, &req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	resp, err := h.client.Geocode(r.Context(), &req)
	if err != nil {
		log.Printf("Error geocoding: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to geocode")
		return
	}

	if err := utils.WriteJSONResponse(w, http.StatusOK, resp); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
