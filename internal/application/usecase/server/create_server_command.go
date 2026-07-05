package server

import (
	"context"

	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/node"
	"github.com/suproxy/backend/internal/domain/server"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type CreateServerCommand struct {
	serverRepo server.Repository
	nodeRepo   node.Repository
	logger     *logger.Logger
}

func NewCreateServerCommand(
	serverRepo server.Repository,
	nodeRepo node.Repository,
	logger *logger.Logger,
) *CreateServerCommand {
	return &CreateServerCommand{
		serverRepo: serverRepo,
		nodeRepo:   nodeRepo,
		logger:     logger,
	}
}

func (c *CreateServerCommand) Execute(ctx context.Context, req *dto.CreateServerRequest) (*dto.ServerResponse, error) {
	// Check if server with same hostname already exists
	existing, _ := c.serverRepo.FindByHostname(ctx, req.Hostname)
	if existing != nil {
		return nil, server.ErrServerAlreadyExists
	}

	// Create server entity
	srv, err := server.NewServer(
		req.Name,
		req.Country,
		req.City,
		req.Hostname,
		req.Provider,
		req.IPv4,
	)
	if err != nil {
		return nil, err
	}

	// Set optional fields
	if req.IPv6 != "" {
		srv.UpdateIPv6(req.IPv6)
	}

	if !req.IsPublic {
		srv.MakePrivate()
	}

	// Save to repository
	if err := c.serverRepo.Create(ctx, srv); err != nil {
		c.logger.Error("Failed to create server", "error", err)
		return nil, err
	}

	c.logger.Info("Server created successfully", "server_id", srv.ID, "hostname", srv.Hostname)

	return mapper.ToServerResponse(srv, 0), nil
}
