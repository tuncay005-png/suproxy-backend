package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/suproxy/backend/internal/infrastructure/jwt"
	"github.com/suproxy/backend/internal/interfaces/http/response"
)

// AuthMiddleware validates JWT token and sets user context
func AuthMiddleware(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader(jwt.AuthorizationHeader)
		if authHeader == "" {
			response.Unauthorized(c, "authorization header is required")
			c.Abort()
			return
		}

		// Validate Bearer prefix
		if !strings.HasPrefix(authHeader, jwt.BearerPrefix) {
			response.Unauthorized(c, "invalid authorization header format. use: Bearer <token>")
			c.Abort()
			return
		}

		// Extract token string
		tokenString := strings.TrimPrefix(authHeader, jwt.BearerPrefix)
		if tokenString == "" {
			response.Unauthorized(c, "token is required")
			c.Abort()
			return
		}

		// Validate access token
		claims, err := jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			handleTokenError(c, err)
			c.Abort()
			return
		}

		// Set user context
		c.Set(jwt.UserIDKey, claims.UserID)
		c.Set(jwt.EmailKey, claims.Email)
		c.Set(jwt.RoleKey, claims.Role)

		c.Next()
	}
}

// RequireAuth ensures user is authenticated
func RequireAuth(jwtManager *jwt.Manager) gin.HandlerFunc {
	return AuthMiddleware(jwtManager)
}

// RequireAdmin ensures user is authenticated and has admin role
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := jwt.GetRole(c)
		if !exists {
			response.Unauthorized(c, "authentication required")
			c.Abort()
			return
		}

		if role != "admin" {
			response.Forbidden(c, "admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth validates token if present but doesn't require it
func OptionalAuth(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(jwt.AuthorizationHeader)
		if authHeader == "" {
			// No auth header, continue without authentication
			c.Next()
			return
		}

		if !strings.HasPrefix(authHeader, jwt.BearerPrefix) {
			// Invalid format, but optional so continue
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, jwt.BearerPrefix)
		claims, err := jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			// Invalid token, but optional so continue
			c.Next()
			return
		}

		// Set user context if valid
		c.Set(jwt.UserIDKey, claims.UserID)
		c.Set(jwt.EmailKey, claims.Email)
		c.Set(jwt.RoleKey, claims.Role)

		c.Next()
	}
}

// handleTokenError returns appropriate error response based on error type
func handleTokenError(c *gin.Context, err error) {
	switch err {
	case jwt.ErrExpiredToken:
		response.ErrorResponse(c, 401, "TOKEN_EXPIRED", "access token has expired")
	case jwt.ErrInvalidToken:
		response.ErrorResponse(c, 401, "INVALID_TOKEN", "invalid access token")
	case jwt.ErrInvalidSignature:
		response.ErrorResponse(c, 401, "INVALID_SIGNATURE", "invalid token signature")
	default:
		response.Unauthorized(c, "token validation failed")
	}
}
