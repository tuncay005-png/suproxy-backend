package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/node"
	"github.com/suproxy/backend/internal/domain/server"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type DeleteServerCommand struct {
	serverRepo server.Repository
	nodeRepo   node.Repository
	logger     *logger.Logger
}

func NewDeleteServerCommand(
	serverRepo server.Repository,
	nodeRepo node.Repository,
	logger *logger.Logger,
) *DeleteServerCommand {
	return &DeleteServerCommand{
		serverRepo: serverRepo,
		nodeRepo:   nodeRepo,
		logger:     logger,
	}
}

func (c *DeleteServerCommand) Execute(ctx context.Context, serverID uuid.UUID) error {
	// Check if server exists
	_, err := c.serverRepo.FindByID(ctx, serverID)
	if err != nil {
		return err
	}

	// Delete server (cascade will delete nodes)
	if err := c.serverRepo.Delete(ctx, serverID); err != nil {
		c.logger.Error("Failed to delete server", "error", err, "server_id", serverID)
		return err
	}

	c.logger.Info("Server deleted successfully", "server_id", serverID)

	return nil
}
