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
	log.Printf("Routing service at %s", cfg.RoutingServiceAddr)
	routingService, err := pygrpc.NewClient(cfg.RoutingServiceAddr)
	if err != nil {
		log.Fatalf("Failed to connect to routing service: %v", err)
	}

	working, err := routingService.HealthCheck(context.Background())
	if err != nil {
		log.Fatalf("Service not working : %t, err : %s", working, err.Error())
	}

	if working {
		log.Printf("Connection to grpc server established successfully at %s", cfg.RoutingServiceAddr)
	}

	defer routingService.Close()

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
