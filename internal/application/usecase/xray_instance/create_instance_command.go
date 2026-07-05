package xray_instance

import (
	"context"

	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/node"
	"github.com/suproxy/backend/internal/domain/xray"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type CreateInstanceCommand struct {
	instanceRepo xray.XrayInstanceRepository
	nodeRepo     node.Repository
	logger       *logger.Logger
}

func NewCreateInstanceCommand(
	instanceRepo xray.XrayInstanceRepository,
	nodeRepo node.Repository,
	logger *logger.Logger,
) *CreateInstanceCommand {
	return &CreateInstanceCommand{
		instanceRepo: instanceRepo,
		nodeRepo:     nodeRepo,
		logger:       logger,
	}
}

func (c *CreateInstanceCommand) Execute(ctx context.Context, req *dto.CreateXrayInstanceRequest) (*dto.XrayInstanceResponse, error) {
	// Verify node exists
	_, err := c.nodeRepo.FindByID(ctx, req.NodeID)
	if err != nil {
		return nil, err
	}

	// Check if instance already exists for this node
	existing, _ := c.instanceRepo.FindByNodeID(ctx, req.NodeID)
	if existing != nil {
		return nil, xray.ErrInstanceAlreadyExists
	}

	// Create instance entity
	instance, err := xray.NewXrayInstance(req.NodeID, req.Version)
	if err != nil {
		return nil, err
	}

	// Save to repository
	if err := c.instanceRepo.Create(ctx, instance); err != nil {
		c.logger.Error("Failed to create xray instance", "error", err)
		return nil, err
	}

	c.logger.Info("Xray instance created successfully", "instance_id", instance.ID, "node_id", req.NodeID)

	return mapper.ToXrayInstanceResponse(instance), nil
}
