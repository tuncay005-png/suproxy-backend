package node

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/node"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type UpdateNodeCommand struct {
	nodeRepo node.Repository
	logger   *logger.Logger
}

func NewUpdateNodeCommand(
	nodeRepo node.Repository,
	logger *logger.Logger,
) *UpdateNodeCommand {
	return &UpdateNodeCommand{
		nodeRepo: nodeRepo,
		logger:   logger,
	}
}

func (c *UpdateNodeCommand) Execute(ctx context.Context, nodeID uuid.UUID, req *dto.UpdateNodeRequest) (*dto.NodeResponse, error) {
	// Find node
	n, err := c.nodeRepo.FindByID(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	// Update max users
	if req.MaxUsers != nil {
		if err := n.UpdateMaxUsers(*req.MaxUsers); err != nil {
			return nil, err
		}
	}

	// Update bandwidth limit
	if req.BandwidthLimitGB != nil {
		if err := n.UpdateBandwidthLimit(*req.BandwidthLimitGB); err != nil {
			return nil, err
		}
	}

	// Update version
	if req.Version != "" {
		n.UpdateVersion(req.Version)
	}

	// Save changes
	if err := c.nodeRepo.Update(ctx, n); err != nil {
		c.logger.Error("Failed to update node", "error", err, "node_id", nodeID)
		return nil, err
	}

	c.logger.Info("Node updated successfully", "node_id", nodeID)

	return mapper.ToNodeResponse(n), nil
}
