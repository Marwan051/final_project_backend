package v1

import (
	"net/http"
	"time"

	"github.com/Marwan051/final_project_backend/internal/api/v1/handlers"
	"github.com/Marwan051/final_project_backend/internal/service/db_tools"
	"github.com/Marwan051/final_project_backend/internal/service/geocoding"
	route_service "github.com/Marwan051/final_project_backend/internal/service/routing"
	"github.com/Marwan051/final_project_backend/internal/service/traffic_updater"
	"github.com/Marwan051/final_project_backend/internal/utils"
)

// NewRouter returns a new router with all v1 API routes
func NewRouter(
	routingService route_service.Router,
	dbToolsService db_tools.DbTools,
	geocodingService geocoding.Geocoding,
	trafficService traffic_updater.TrafficUpdater,
) *http.ServeMux {
	mux := http.NewServeMux()

	routingHandler := handlers.NewRoutingHandler(routingService)
	dbToolsHandler := handlers.NewDbToolsHandler(dbToolsService)
	geocodingHandler := handlers.NewGeocodingHandler(geocodingService)
	trafficHandler := handlers.NewTrafficHandler(trafficService)

	mux.HandleFunc("GET /health", HealthHandler)

	// Routing endpoint
	mux.HandleFunc("POST /route", routingHandler.FindRoute)

	// Geocoding endpoints
	mux.HandleFunc("POST /geocode", geocodingHandler.Geocode)

	// DB Tools endpoints
	mux.HandleFunc("POST /nearby-trips", dbToolsHandler.NearbyTrips)

	// Traffic Updater endpoints
	mux.HandleFunc("POST /traffic/trigger", trafficHandler.TriggerUpdate)
	mux.HandleFunc("GET /traffic/status", trafficHandler.GetStatus)
	mux.HandleFunc("POST /traffic/update-trip", trafficHandler.UpdateTrip)
	mux.HandleFunc("POST /traffic/street", trafficHandler.StreetTraffic)
	mux.HandleFunc("GET /traffic/streets", trafficHandler.ListStreets)

	return mux
}

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// HealthHandler returns the health status of the service
// @Summary      Health Check
// @Description  Returns the health status of the backend service
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /health [get]
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}
