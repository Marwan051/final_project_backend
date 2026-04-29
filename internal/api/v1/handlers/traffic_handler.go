package handlers

import (
	"log"
	"net/http"

	"github.com/Marwan051/final_project_backend/internal/service/traffic_updater"
	pb "github.com/Marwan051/final_project_backend/internal/service/traffic_updater/proto"
	"github.com/Marwan051/final_project_backend/internal/utils"
)

type TrafficHandler struct {
	client traffic_updater.TrafficUpdater
}

func NewTrafficHandler(client traffic_updater.TrafficUpdater) *TrafficHandler {
	return &TrafficHandler{
		client: client,
	}
}

// TriggerUpdate starts a traffic update manually
// @Summary      Trigger Traffic Update
// @Description  Triggers a background traffic update process
// @Tags         traffic
// @Accept       json
// @Produce      json
// @Param        request body pb.TriggerRequest true "Trigger Request"
// @Success      200  {object}  pb.UpdateResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/v1/traffic/trigger [post]
func (h *TrafficHandler) TriggerUpdate(w http.ResponseWriter, r *http.Request) {
	var req pb.TriggerRequest
	if err := utils.DecodeJSONBody(r, &req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	resp, err := h.client.TriggerUpdate(r.Context(), &req)
	if err != nil {
		log.Printf("Error triggering update: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to trigger update")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, resp)
}

// GetStatus returns the current status of traffic services
// @Summary      Get Traffic Status
// @Description  Retrieves the current traffic update status
// @Tags         traffic
// @Produce      json
// @Success      200  {object}  pb.StatusResponse
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/v1/traffic/status [get]
func (h *TrafficHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	resp, err := h.client.GetStatus(r.Context())
	if err != nil {
		log.Printf("Error getting status: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to get status")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, resp)
}

// UpdateTrip updates the traffic impact based on a specific trip
// @Summary      Update Trip Traffic
// @Description  Submits trip information to affect current traffic modelling
// @Tags         traffic
// @Accept       json
// @Produce      json
// @Param        request body pb.UpdateTripRequest true "Update Trip Request"
// @Success      200  {object}  pb.UpdateResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/v1/traffic/update-trip [post]
func (h *TrafficHandler) UpdateTrip(w http.ResponseWriter, r *http.Request) {
	var req pb.UpdateTripRequest
	if err := utils.DecodeJSONBody(r, &req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	resp, err := h.client.UpdateTrip(r.Context(), &req)
	if err != nil {
		log.Printf("Error updating trip: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to update trip")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, resp)
}

// StreetTraffic gets traffic for a specific street
// @Summary      Get Street Traffic
// @Description  Fetches the current calculated traffic load on a specific street
// @Tags         traffic
// @Accept       json
// @Produce      json
// @Param        request body pb.StreetTrafficRequest true "Get Street Traffic Request"
// @Success      200  {object}  pb.StreetTrafficResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/v1/traffic/street [post]
func (h *TrafficHandler) StreetTraffic(w http.ResponseWriter, r *http.Request) {
	var req pb.StreetTrafficRequest
	if err := utils.DecodeJSONBody(r, &req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	resp, err := h.client.StreetTraffic(r.Context(), &req)
	if err != nil {
		log.Printf("Error getting street traffic: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to get street traffic")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, resp)
}

// ListStreets retrieves all tracked streets and their traffic states
// @Summary      List All Streets Traffic
// @Description  Returns a list of all streets actively tracked and their traffic levels
// @Tags         traffic
// @Produce      json
// @Success      200  {object}  pb.StreetListResponse
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/v1/traffic/streets [get]
func (h *TrafficHandler) ListStreets(w http.ResponseWriter, r *http.Request) {
	resp, err := h.client.ListStreets(r.Context())
	if err != nil {
		log.Printf("Error listing streets: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to list streets")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, resp)
}
