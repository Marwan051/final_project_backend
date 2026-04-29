package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Marwan051/final_project_backend/internal/auth"
	"github.com/Marwan051/final_project_backend/internal/server"
	db_tools_client "github.com/Marwan051/final_project_backend/internal/service/db_tools/client"
	geocoding_client "github.com/Marwan051/final_project_backend/internal/service/geocoding/client"
	routing_client "github.com/Marwan051/final_project_backend/internal/service/routing/client"
	traffic_client "github.com/Marwan051/final_project_backend/internal/service/traffic_updater/client"
	"github.com/Marwan051/final_project_backend/internal/utils"
)

func main() {
	// Load configuration
	if err := utils.LoadENV(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	cfg := utils.Cfg

	// load routing service
	routingService, err := routing_client.NewRoutingClient(routing_client.ClientConfig{
		Address: cfg.RoutingServiceAddr,
	})
	if err != nil {
		log.Fatalf("Failed to connect to routing service: %v", err)
	}
	defer routingService.Close()

	// load db tools service
	dbToolsService, err := db_tools_client.NewDbToolsClient(db_tools_client.ClientConfig{
		Address: cfg.DbToolsAddr,
	})
	if err != nil {
		log.Fatalf("Failed to connect to db tools service: %v", err)
	}
	defer dbToolsService.Close()

	// load geocoding service
	geocodingService, err := geocoding_client.NewGeocodingClient(geocoding_client.ClientConfig{
		Address: cfg.GeocodingAddr,
	})
	if err != nil {
		log.Fatalf("Failed to connect to geocoding service: %v", err)
	}
	defer geocodingService.Close()

	// load traffic updater service
	trafficService, err := traffic_client.NewTrafficUpdaterClient(traffic_client.ClientConfig{
		Address: cfg.TrafficUpdaterAddr,
	})
	if err != nil {
		log.Fatalf("Failed to connect to traffic service: %v", err)
	}
	defer trafficService.Close()

	log.Printf("Waiting for gRPC services to be ready...")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	working, err := routingService.HealthCheck(ctx)
	if err != nil || !working {
		log.Fatalf("gRPC routing health check failed: %v", err)
	}
	working, err = dbToolsService.HealthCheck(ctx)
	if err != nil || !working {
		log.Fatalf("gRPC db tools health check failed: %v", err)
	}
	working, err = geocodingService.HealthCheck(ctx)
	if err != nil || !working {
		log.Fatalf("gRPC geocoding health check failed: %v", err)
	}

	log.Printf("All gRPC connections verified")

	// Initialize Firebase Auth verifier using ADC and optional explicit project ID.
	fbVerifier, err := auth.NewFirebaseVerifierWithProjectID(context.Background(), cfg.EffectiveFirebaseProjectID())
	if err != nil {
		log.Fatalf("failed to initialize firebase auth: %v", err)
	}

	// Create HTTP handler with injected dependencies
	handler := server.NewHandler(routingService, dbToolsService, geocodingService, trafficService, fbVerifier)

	// Create server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on :%s", cfg.Port)
		log.Printf("Server starting in %s mode", cfg.ENV)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
