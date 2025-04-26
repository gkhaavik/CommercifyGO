package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/infrastructure/database"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
	"github.com/zenfulcode/commercify/internal/interfaces/api"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		// Using standard log since logger isn't initialized yet
		log.Println("No .env file found, using environment variables")
	}

	// Initialize logger
	logger := logger.NewLogger()
	logger.Info("Starting Commercify backend service")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.NewPostgresConnection(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run database migrations
	if err := database.RunMigrations(db, cfg.Database); err != nil {
		logger.Fatal("Failed to run database migrations: %v", err)
	}

	// Initialize API server
	server := api.NewServer(cfg, db, logger)

	// Start server in a goroutine
	go func() {
		logger.Info("Starting server on port %s", cfg.Server.Port)
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited properly")
}
