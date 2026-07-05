package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/suproxy/backend/internal/infrastructure/config"
	"github.com/suproxy/backend/internal/infrastructure/logger"
	"github.com/suproxy/backend/internal/infrastructure/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration: %v", err))
	}

	// Initialize logger
	log := logger.New(cfg.Log.Level, cfg.Log.Format)
	defer log.Sync()

	log.Info("Starting SuProxy Backend API",
		"version", "1.0.0",
		"environment", cfg.Environment,
	)

	// Initialize HTTP server
	srv := server.New(cfg, log)

	// Start server
	go func() {
		log.Info("Starting HTTP server", "address", cfg.Server.Address)
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", "error", err)
	}

	log.Info("Server stopped gracefully")
}
