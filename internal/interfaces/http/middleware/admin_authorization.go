package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/suproxy/backend/internal/infrastructure/jwt"
	"github.com/suproxy/backend/internal/infrastructure/logger"
	"github.com/suproxy/backend/internal/interfaces/http/response"
)

// AdminAuthorization ensures user has admin role
// This middleware MUST be used after AuthMiddleware
func AdminAuthorization(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get role from context (set by AuthMiddleware)
		role, exists := jwt.GetRole(c)
		if !exists {
			log.Warn("Admin access attempted without role in context",
				"ip", c.ClientIP(),
				"user_agent", c.GetHeader("User-Agent"),
				"path", c.Request.URL.Path)

			response.Unauthorized(c, "authentication required")
			c.Abort()
			return
		}

		// Check if user is admin
		if role != "admin" {
			// Get user ID for logging
			userID, _ := jwt.GetUserID(c)
			email, _ := jwt.GetEmail(c)

			log.Warn("Non-admin user attempted to access admin endpoint",
				"user_id", userID,
				"email", email,
				"role", role,
				"ip", c.ClientIP(),
				"path", c.Request.URL.Path,
				"method", c.Request.Method)

			response.Forbidden(c, "admin access required")
			c.Abort()
			return
		}

		// User is admin, proceed
		c.Next()
	}
}
