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

func (h *TrafficHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	resp, err := h.client.GetStatus(r.Context())
	if err != nil {
		log.Printf("Error getting status: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to get status")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, resp)
}

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

func (h *TrafficHandler) ListStreets(w http.ResponseWriter, r *http.Request) {
	resp, err := h.client.ListStreets(r.Context())
	if err != nil {
		log.Printf("Error listing streets: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to list streets")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, resp)
}
