package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger is a middleware that logs HTTP requests
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Log after request completes
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		log.Printf("[%s] %s %s | %d | %v | %s",
			method,
			path,
			c.Request.Proto,
			statusCode,
			latency,
			clientIP,
		)
	}
}
