package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/server"
	"github.com/sirupsen/logrus"
)

func main() {
	log.Println("Starting Cloud Consulting Backend...")
	
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	
	log.Printf("Configuration loaded. Port: %s, LogLevel: %d", cfg.Port, cfg.LogLevel)

	// Initialize logger
	logger := logrus.New()
	logger.SetLevel(logrus.Level(cfg.LogLevel))
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Create server
	srv, err := server.New(cfg, logger)
	if err != nil {
		logger.Fatalf("Failed to create server: %v", err)
	}
	
	log.Println("Server created successfully")

	// Start server
	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: srv.Handler(),
	}

	// Graceful shutdown
	go func() {
		log.Printf("Starting HTTP server on port %s", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	logger.Infof("Server started on port %s", cfg.Port)
	log.Printf("Server is running at http://localhost:%s", cfg.Port)
	log.Println("Health check available at: http://localhost:" + cfg.Port + "/health")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}