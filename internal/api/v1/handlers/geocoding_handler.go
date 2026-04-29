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

// Geocode converts an address or place into coordinates
// @Summary      Geocode Address
// @Description  Converts a text address into latitude and longitude coordinates
// @Tags         geocoding
// @Accept       json
// @Produce      json
// @Param        request body pb.GeocodeRequest true "Geocode Request (Address)"
// @Success      200  {object}  pb.GeocodeResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/v1/geocode [post]
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
