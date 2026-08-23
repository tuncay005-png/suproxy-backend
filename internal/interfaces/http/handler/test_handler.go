package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/suproxy/backend/internal/interfaces/http/response"
)

// TestHandler handles test-only endpoints for E2E testing
type TestHandler struct{}

// NewTestHandler creates a new test handler
func NewTestHandler() *TestHandler {
	return &TestHandler{}
}

// SimulateExpiredToken godoc
// @Summary Simulate expired token error
// @Description Returns a TOKEN_EXPIRED error for E2E testing
// @Tags test
// @Accept json
// @Produce json
// @Failure 401 {object} response.Response
// @Router /api/v1/test/expired-token [get]
func (h *TestHandler) SimulateExpiredToken(c *gin.Context) {
	response.ErrorResponse(c, 401, "TOKEN_EXPIRED", "access token has expired (test simulation)")
}
