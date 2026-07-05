package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/infrastructure/jwt"
	"github.com/suproxy/backend/internal/interfaces/http/response"
)

// XrayHandler handles xray-related endpoints
type XrayHandler struct {
	// Use cases will be added when needed
}

// NewXrayHandler creates a new xray handler
func NewXrayHandler() *XrayHandler {
	return &XrayHandler{}
}

// ListInstances godoc
// @Summary List Xray instances
// @Description Get a list of Xray instances for the authenticated user
// @Tags xray
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.XrayInstanceListResponse
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/xray/instances [get]
func (h *XrayHandler) ListInstances(c *gin.Context) {
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

	// TODO: Call use case to list xray instances
	// For now, return empty list
	resp := dto.XrayInstanceListResponse{
		Instances: []*dto.XrayInstanceResponse{},
		Total:     0,
	}

	response.SuccessOK(c, resp)
}
