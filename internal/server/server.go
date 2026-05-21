package server

import (
	"net/http"

	_ "github.com/Marwan051/final_project_backend/docs"
	v1 "github.com/Marwan051/final_project_backend/internal/api/v1"
	authpkg "github.com/Marwan051/final_project_backend/internal/auth"
	agent_service "github.com/Marwan051/final_project_backend/internal/service/agent"
	"github.com/Marwan051/final_project_backend/internal/service/db_tools"
	"github.com/Marwan051/final_project_backend/internal/service/geocoding"
	route_service "github.com/Marwan051/final_project_backend/internal/service/routing"
	"github.com/Marwan051/final_project_backend/internal/service/traffic_updater"
	"github.com/Marwan051/final_project_backend/internal/utils"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// NewHandler creates the application's HTTP handler with middleware
func NewHandler(
	routingService route_service.Router,
	agentService agent_service.Agent,
	dbToolsService db_tools.DbTools,
	geocodingService geocoding.Geocoding,
	trafficService traffic_updater.TrafficUpdater,
	verifier authpkg.Verifier,
) http.Handler {
	// Create v1 router with dependencies
	v1Router := v1.NewRouter(routingService, agentService, dbToolsService, geocodingService, trafficService)

	// Apply auth specifically to the v1 sub-router
	protectedV1Router := Auth(verifier)(v1Router)

	// Main router
	mux := http.NewServeMux()

	// Unauthenticated routes
	mux.HandleFunc("GET /health", v1.HealthHandler)

	// Swagger UI endpoint (unauthenticated).
	// Enabled when running in dev *or* when `ENABLE_SWAGGER=true` is set.
	if utils.Cfg.ENV == "dev" || utils.Cfg.EnableSwagger {
		mux.Handle("/docs/", httpSwagger.Handler(
			httpSwagger.URL("/docs/doc.json"), //The url pointing to API definition
			httpSwagger.UIConfig(map[string]string{
				"persistAuthorization": "true",
				"requestInterceptor": `function(req) {
  if (req.headers && req.headers.Authorization && !req.headers.Authorization.startsWith('Bearer ')) {
    req.headers.Authorization = 'Bearer ' + req.headers.Authorization;
  }
  return req;
}`,
			}),
		))
	}

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
