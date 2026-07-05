package node

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/node"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type UpdateNodeMetricsCommand struct {
	nodeRepo node.Repository
	logger   *logger.Logger
}

func NewUpdateNodeMetricsCommand(
	nodeRepo node.Repository,
	logger *logger.Logger,
) *UpdateNodeMetricsCommand {
	return &UpdateNodeMetricsCommand{
		nodeRepo: nodeRepo,
		logger:   logger,
	}
}

func (c *UpdateNodeMetricsCommand) Execute(ctx context.Context, nodeID uuid.UUID, req *dto.UpdateNodeMetricsRequest) (*dto.NodeResponse, error) {
	// Find node
	n, err := c.nodeRepo.FindByID(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	// Update metrics
	if err := n.UpdateMetrics(req.CPUUsage, req.RAMUsage, req.LatencyMs); err != nil {
		return nil, err
	}

	// Save changes
	if err := c.nodeRepo.Update(ctx, n); err != nil {
		c.logger.Error("Failed to update node metrics", "error", err, "node_id", nodeID)
		return nil, err
	}

	c.logger.Debug("Node metrics updated", "node_id", nodeID, "health", n.HealthStatus)

	return mapper.ToNodeResponse(n), nil
}
