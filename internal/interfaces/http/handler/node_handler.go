package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/suproxy/backend/internal/application/usecase/node"
	"github.com/suproxy/backend/internal/interfaces/http/response"
)

// NodeHandler handles node-related endpoints
type NodeHandler struct {
	listNodesQuery *node.ListNodesQuery
}

// NewNodeHandler creates a new node handler
func NewNodeHandler(
	listNodesQuery *node.ListNodesQuery,
) *NodeHandler {
	return &NodeHandler{
		listNodesQuery: listNodesQuery,
	}
}

// ListNodes godoc
// @Summary List all nodes
// @Description Get a list of all available VPN nodes
// @Tags nodes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {object} dto.NodeListResponse
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/nodes [get]
func (h *NodeHandler) ListNodes(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	nodes, err := h.listNodesQuery.Execute(c.Request.Context(), offset, limit)
	if err != nil {
		response.InternalError(c, "failed to fetch nodes")
		return
	}

	response.SuccessOK(c, nodes)
}
