package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/node"
	"github.com/suproxy/backend/internal/domain/server"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type UpdateServerCommand struct {
	serverRepo server.Repository
	nodeRepo   node.Repository
	logger     *logger.Logger
}

func NewUpdateServerCommand(
	serverRepo server.Repository,
	nodeRepo node.Repository,
	logger *logger.Logger,
) *UpdateServerCommand {
	return &UpdateServerCommand{
		serverRepo: serverRepo,
		nodeRepo:   nodeRepo,
		logger:     logger,
	}
}

func (c *UpdateServerCommand) Execute(ctx context.Context, serverID uuid.UUID, req *dto.UpdateServerRequest) (*dto.ServerResponse, error) {
	// Find server
	srv, err := c.serverRepo.FindByID(ctx, serverID)
	if err != nil {
		return nil, err
	}

	// Update details
	if req.Name != "" || req.City != "" || req.Provider != "" {
		srv.UpdateDetails(req.Name, req.City, req.Provider)
	}

	// Update IPv6
	if req.IPv6 != "" {
		srv.UpdateIPv6(req.IPv6)
	}

	// Update public status
	if req.IsPublic != nil {
		if *req.IsPublic {
			srv.MakePublic()
		} else {
			srv.MakePrivate()
		}
	}

	// Save changes
	if err := c.serverRepo.Update(ctx, srv); err != nil {
		c.logger.Error("Failed to update server", "error", err, "server_id", serverID)
		return nil, err
	}

	// Get node count
	nodeCount, _ := c.nodeRepo.CountByServerID(ctx, serverID)

	c.logger.Info("Server updated successfully", "server_id", serverID)

	return mapper.ToServerResponse(srv, int(nodeCount)), nil
}
