package node

import (
	"context"

	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/node"
	"github.com/suproxy/backend/internal/domain/server"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type CreateNodeCommand struct {
	nodeRepo   node.Repository
	serverRepo server.Repository
	logger     *logger.Logger
}

func NewCreateNodeCommand(
	nodeRepo node.Repository,
	serverRepo server.Repository,
	logger *logger.Logger,
) *CreateNodeCommand {
	return &CreateNodeCommand{
		nodeRepo:   nodeRepo,
		serverRepo: serverRepo,
		logger:     logger,
	}
}

func (c *CreateNodeCommand) Execute(ctx context.Context, req *dto.CreateNodeRequest) (*dto.NodeResponse, error) {
	// Verify server exists
	srv, err := c.serverRepo.FindByID(ctx, req.ServerID)
	if err != nil {
		return nil, err
	}

	// Check if server is available
	if !srv.IsAvailable() {
		return nil, server.ErrServerNotAvailable
	}

	// Create node entity
	n, err := node.NewNode(
		req.ServerID,
		node.Protocol(req.Protocol),
		req.Port,
		req.MaxUsers,
		req.BandwidthLimitGB,
	)
	if err != nil {
		return nil, err
	}

	// Save to repository
	if err := c.nodeRepo.Create(ctx, n); err != nil {
		c.logger.Error("Failed to create node", "error", err)
		return nil, err
	}

	c.logger.Info("Node created successfully", "node_id", n.ID, "server_id", req.ServerID)

	return mapper.ToNodeResponse(n), nil
}
