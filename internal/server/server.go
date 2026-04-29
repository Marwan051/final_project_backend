package server

import (
	"net/http"

	v1 "github.com/Marwan051/final_project_backend/internal/api/v1"
	authpkg "github.com/Marwan051/final_project_backend/internal/auth"
	"github.com/Marwan051/final_project_backend/internal/service/db_tools"
	"github.com/Marwan051/final_project_backend/internal/service/geocoding"
	route_service "github.com/Marwan051/final_project_backend/internal/service/routing"
	"github.com/Marwan051/final_project_backend/internal/service/traffic_updater"
)

// NewHandler creates the application's HTTP handler with middleware
func NewHandler(
	routingService route_service.Router,
	dbToolsService db_tools.DbTools,
	geocodingService geocoding.Geocoding,
	trafficService traffic_updater.TrafficUpdater,
	verifier authpkg.Verifier,
) http.Handler {
	// Create v1 router with dependencies
	v1Router := v1.NewRouter(routingService, dbToolsService, geocodingService, trafficService)

	// Apply auth specifically to the v1 sub-router
	protectedV1Router := Auth(verifier)(v1Router)

	// Main router
	mux := http.NewServeMux()
	
	// Unauthenticated routes
	mux.HandleFunc("GET /health", v1.HealthHandler)

	// Authenticated routes
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", protectedV1Router))

	// Apply remaining middleware globally (Headers first, then PanicRecover, Logging)
	handler := ChainMiddleware(mux,
		Headers,
		PanicRecover,
		Logging,
	)

	return handler
}