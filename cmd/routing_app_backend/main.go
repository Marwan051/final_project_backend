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
	"github.com/Marwan051/final_project_backend/internal/utils"
)

func main() {
	// Load configuration
	if err := utils.LoadENV(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	cfg := utils.Cfg

	// Create handler from routes
	handler := server.NewHandler()

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
