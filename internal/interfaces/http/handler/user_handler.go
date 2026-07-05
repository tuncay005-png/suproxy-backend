package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/infrastructure/jwt"
	"github.com/suproxy/backend/internal/interfaces/http/response"
)

// UserHandler handles user-related endpoints
type UserHandler struct {
	// Use cases will be added when implemented
}

// NewUserHandler creates a new user handler
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetMe godoc
// @Summary Get current user profile
// @Description Get the profile of the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} response.Response
// @Router /api/v1/users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userIDStr, exists := jwt.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "user not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	// TODO: Call use case to get user
	// For now, return a placeholder response
	resp := dto.UserResponse{
		ID:        userID,
		Email:     "user@example.com",
		CreatedAt: time.Now(),
	}

	response.SuccessOK(c, resp)
}

// UpdateMe godoc
// @Summary Update current user profile
// @Description Update the profile of the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateUserRequest true "Update user request"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/users/me [put]
func (h *UserHandler) UpdateMe(c *gin.Context) {
	userIDStr, exists := jwt.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "user not authenticated")
		return
	}

	_, err := uuid.Parse(userIDStr)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// TODO: Call use case to update user
	response.SuccessOK(c, map[string]string{
		"message": "user updated successfully",
	})
}
