package server

import (
	"net/http"

	v1 "github.com/Marwan051/final_project_backend/internal/api/v1"
	authpkg "github.com/Marwan051/final_project_backend/internal/auth"
	"github.com/Marwan051/final_project_backend/internal/service/route_service"
)

// NewHandler creates the application's HTTP handler with middleware
func NewHandler(routingService route_service.Router, verifier authpkg.Verifier) http.Handler {
	// Create v1 router with dependencies
	v1Router := v1.NewRouter(routingService)

	// Main router
	mux := http.NewServeMux()
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1Router))

	// Apply middleware (Headers first, then Auth globally)
	handler := ChainMiddleware(mux,
		Headers,
		Auth(verifier),
		PanicRecover,
		Logging,
	)

	return handler
}
