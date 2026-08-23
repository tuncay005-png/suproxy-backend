package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS returns a middleware that handles CORS with environment-based origin validation
func CORS() gin.HandlerFunc {
	// Parse allowed origins from environment variable
	// Format: CORS_ALLOWED_ORIGINS=https://admin.example.com,https://app.example.com
	allowedOriginsEnv := os.Getenv("CORS_ALLOWED_ORIGINS")

	// Default to localhost for development if not specified
	if allowedOriginsEnv == "" {
		allowedOriginsEnv = "http://localhost:3000,http://localhost:3001"
	}

	// Parse into slice
	allowedOrigins := strings.Split(allowedOriginsEnv, ",")

	// Trim whitespace from each origin
	for i := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		// If origin is allowed, set CORS headers
		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		}

		// Handle preflight OPTIONS request
		if c.Request.Method == "OPTIONS" {
			if allowed {
				c.AbortWithStatus(204)
			} else {
				c.AbortWithStatus(403)
			}
			return
		}

		c.Next()
	}
}
