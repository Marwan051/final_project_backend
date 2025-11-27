package server

import (
	"net/http"

	v1 "github.com/Marwan051/final_project_backend/internal/api/v1"
)

// NewHandler creates the application's HTTP handler with middleware
func NewHandler() http.Handler {
	mux := http.NewServeMux()

	// Mount API v1 routes under /api/v1/
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1.NewRouter()))

	// Apply middleware chain (order: first runs first)
	return ChainMiddleware(mux,
		Headers,
		PanicRecover,
		Logging,
	)
}
