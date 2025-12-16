package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Marwan051/final_project_backend/internal/server"
	"github.com/Marwan051/final_project_backend/internal/service/route_service/pygrpc"
	"github.com/Marwan051/final_project_backend/internal/utils"
)

func main() {
	// Load configuration
	if err := utils.LoadENV(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	cfg := utils.Cfg
	routingService, err := pygrpc.NewClient(pygrpc.ClientConfig{
		Address: cfg.RoutingServiceAddr,
	})
	if err != nil {
		log.Fatalf("Failed to connect to routing service: %v", err)
	}
	defer routingService.Close()

	log.Printf("Waiting for gRPC service to be ready...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	working, err := routingService.HealthCheck(ctx)
	if err != nil {
		log.Fatalf("gRPC health check failed: %v", err)
	}
	if !working {
		log.Fatalf("gRPC service reported unhealthy")
	}
	log.Printf("gRPC connection verified at %s", cfg.RoutingServiceAddr)

	// Create HTTP handler with injected dependencies
	handler := server.NewHandler(routingService)

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
