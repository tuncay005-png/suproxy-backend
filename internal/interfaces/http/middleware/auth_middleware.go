package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/suproxy/backend/internal/infrastructure/jwt"
	"github.com/suproxy/backend/internal/interfaces/http/response"
)

func AuthMiddleware(jwtManager *jwt.Manager) gin.HandlerFunc {
	return jwt.AuthMiddleware(jwtManager)
}

func RequireAuth(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := jwt.GetUserID(c)
		if !exists {
			response.Unauthorized(c, "authentication required")
			c.Abort()
			return
		}
		c.Next()
	}
}
