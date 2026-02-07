package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aaronbengochea/periscope/backend-go/config"
	"github.com/aaronbengochea/periscope/backend-go/internal/api"
	"github.com/aaronbengochea/periscope/backend-go/pkg/database"
	"github.com/aaronbengochea/periscope/backend-go/pkg/massive"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// Initialize Massive API client
	massiveClient := massive.NewClient(cfg.MassiveBaseURL, cfg.MassiveAPIKey)
	log.Println("✓ Initialized Massive API client")

	// TODO: Initialize database connection
	// For now, we'll skip database connection since we need the actual connection string
	// This will be implemented in the next step
	var db *database.DB
	if cfg.DatabaseURL != "" {
		db, err = database.NewSupabaseDB(cfg.DatabaseURL)
		if err != nil {
			log.Printf("Warning: Failed to connect to database: %v", err)
			log.Println("Continuing without database connection...")
		} else {
			defer db.Close()
			log.Println("✓ Connected to Supabase database")
		}
	}

	// Setup router
	router := api.NewRouter(cfg, db, massiveClient)

	// Create HTTP server
	addr := fmt.Sprintf(":%s", cfg.Port)
	srv := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Server starting on http://localhost%s", addr)
		log.Printf("   Environment: %s", cfg.GinMode)
		log.Printf("   Massive API: %s", cfg.MassiveBaseURL)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
