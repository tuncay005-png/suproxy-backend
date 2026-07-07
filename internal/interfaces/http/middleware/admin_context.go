package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/infrastructure/jwt"
)

// AdminContext provides helper functions for admin operations
type AdminContext struct {
	UserID    uuid.UUID
	Email     string
	Role      string
	IPAddress string
	UserAgent string
}

// GetAdminContext extracts admin context from gin context
func GetAdminContext(c *gin.Context) (*AdminContext, error) {
	// Get user ID
	userIDStr, exists := jwt.GetUserID(c)
	if !exists {
		return nil, ErrMissingUserID
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ErrInvalidUserID
	}

	// Get email
	email, exists := jwt.GetEmail(c)
	if !exists {
		return nil, ErrMissingEmail
	}

	// Get role
	role, exists := jwt.GetRole(c)
	if !exists {
		return nil, ErrMissingRole
	}

	// Get IP address
	ipAddress := c.ClientIP()

	// Get user agent
	userAgent := c.GetHeader("User-Agent")

	return &AdminContext{
		UserID:    userID,
		Email:     email,
		Role:      role,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}, nil
}

// IsAdmin checks if the context user is an admin
func (ac *AdminContext) IsAdmin() bool {
	return ac.Role == "admin"
}
