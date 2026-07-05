package node

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/node"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type DeleteNodeCommand struct {
	nodeRepo node.Repository
	logger   *logger.Logger
}

func NewDeleteNodeCommand(
	nodeRepo node.Repository,
	logger *logger.Logger,
) *DeleteNodeCommand {
	return &DeleteNodeCommand{
		nodeRepo: nodeRepo,
		logger:   logger,
	}
}

func (c *DeleteNodeCommand) Execute(ctx context.Context, nodeID uuid.UUID) error {
	// Check if node exists
	_, err := c.nodeRepo.FindByID(ctx, nodeID)
	if err != nil {
		return err
	}

	// Delete node
	if err := c.nodeRepo.Delete(ctx, nodeID); err != nil {
		c.logger.Error("Failed to delete node", "error", err, "node_id", nodeID)
		return err
	}

	c.logger.Info("Node deleted successfully", "node_id", nodeID)

	return nil
}
