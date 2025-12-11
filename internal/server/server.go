package server

import (
	"net/http"

	v1 "github.com/Marwan051/final_project_backend/internal/api/v1"
	"github.com/Marwan051/final_project_backend/internal/service/route_service"
)

// NewHandler creates the application's HTTP handler with middleware
func NewHandler(routingService route_service.Router) http.Handler {
	// Create v1 router with dependencies
	v1Router := v1.NewRouter(routingService)

	// Main router
	mux := http.NewServeMux()
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1Router))

	// Apply middleware
	handler := ChainMiddleware(mux,
		Headers,
		PanicRecover,
		Logging,
	)

	return handler
}
