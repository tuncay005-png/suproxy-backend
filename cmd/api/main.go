package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/suproxy/backend/internal/infrastructure/bootstrap"
	"github.com/suproxy/backend/internal/infrastructure/server"
)

func main() {
	// Initialize application
	app, err := bootstrap.Initialize()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize application: %v", err))
	}
	defer app.Shutdown()

	app.Logger.Info("Starting SuProxy Backend API",
		"version", "1.0.0",
		"environment", app.Config.Environment,
	)

	// Initialize HTTP server
	srv := server.New(app)

	// Start server
	go func() {
		app.Logger.Info("Starting HTTP server", "address", app.Config.Server.Address)
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			app.Logger.Fatal("Failed to start server", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	app.Logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		app.Logger.Fatal("Server forced to shutdown", "error", err)
	}

	app.Logger.Info("Server stopped gracefully")
}
