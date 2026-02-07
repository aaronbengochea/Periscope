package api

import (
	"net/http"

	"github.com/aaronbengochea/periscope/backend-go/config"
	"github.com/aaronbengochea/periscope/backend-go/internal/api/handlers"
	"github.com/aaronbengochea/periscope/backend-go/internal/api/middleware"
	"github.com/aaronbengochea/periscope/backend-go/pkg/database"
	"github.com/aaronbengochea/periscope/backend-go/pkg/massive"
	"github.com/gin-gonic/gin"
)

// NewRouter creates and configures the HTTP router
func NewRouter(cfg *config.Config, db *database.DB, massiveClient *massive.Client) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())                     // Recover from panics
	router.Use(middleware.Logger())                // Structured logging
	router.Use(middleware.CORS())                  // CORS for frontend

	// Health check endpoint (supports both GET and HEAD for Docker healthcheck)
	healthHandler := func(c *gin.Context) {
		status := gin.H{
			"status": "healthy",
			"service": "periscope-api",
		}

		// Check database connection if available
		if db != nil {
			if err := db.Health(c.Request.Context()); err != nil {
				status["database"] = "unhealthy"
				status["database_error"] = err.Error()
				c.JSON(http.StatusServiceUnavailable, status)
				return
			}
			status["database"] = "healthy"
		} else {
			status["database"] = "not_connected"
		}

		c.JSON(http.StatusOK, status)
	}
	router.GET("/health", healthHandler)
	router.HEAD("/health", healthHandler)

	// Initialize handlers
	optionsHandler := handlers.NewOptionsHandler(massiveClient)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Options endpoints
		v1.GET("/options/:ticker", optionsHandler.GetOptionsChain)

		// Portfolio endpoints (to be implemented)
		v1.GET("/portfolio", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Portfolio endpoint - coming soon",
			})
		})
	}

	return router
}
