package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/suproxy/backend/internal/application/usecase/server"
	"github.com/suproxy/backend/internal/interfaces/http/response"
)

// ServerHandler handles server-related endpoints
type ServerHandler struct {
	listServersQuery *server.ListServersQuery
}

// NewServerHandler creates a new server handler
func NewServerHandler(
	listServersQuery *server.ListServersQuery,
) *ServerHandler {
	return &ServerHandler{
		listServersQuery: listServersQuery,
	}
}

// ListServers godoc
// @Summary List all servers
// @Description Get a list of all available VPN servers
// @Tags servers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {object} dto.ServerListResponse
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/servers [get]
func (h *ServerHandler) ListServers(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	servers, err := h.listServersQuery.Execute(c.Request.Context(), offset, limit)
	if err != nil {
		response.InternalError(c, "failed to fetch servers")
		return
	}

	response.SuccessOK(c, servers)
}
