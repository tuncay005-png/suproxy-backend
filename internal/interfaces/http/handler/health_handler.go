package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/suproxy/backend/internal/infrastructure/database"
	"github.com/suproxy/backend/internal/interfaces/http/response"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	db *database.Database
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *database.Database) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
}

// Health godoc
// @Summary Health check
// @Description Check if the API is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	dbStatus := "healthy"
	if err := h.db.Ping(); err != nil {
		dbStatus = "unhealthy"
	}

	healthResp := HealthResponse{
		Status: "ok",
		Services: map[string]string{
			"database": dbStatus,
		},
	}

	response.SuccessOK(c, healthResp)
}

// Ready godoc
// @Summary Readiness check
// @Description Check if the API is ready to accept requests
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /ready [get]
func (h *HealthHandler) Ready(c *gin.Context) {
	if err := h.db.Ping(); err != nil {
		response.InternalError(c, "database not ready")
		return
	}

	healthResp := HealthResponse{
		Status: "ready",
		Services: map[string]string{
			"database": "connected",
		},
	}

	response.SuccessOK(c, healthResp)
}
