package inbound

import (
	"context"

	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/xray"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type CreateInboundCommand struct {
	inboundRepo  xray.InboundRepository
	instanceRepo xray.XrayInstanceRepository
	logger       *logger.Logger
}

func NewCreateInboundCommand(
	inboundRepo xray.InboundRepository,
	instanceRepo xray.XrayInstanceRepository,
	logger *logger.Logger,
) *CreateInboundCommand {
	return &CreateInboundCommand{
		inboundRepo:  inboundRepo,
		instanceRepo: instanceRepo,
		logger:       logger,
	}
}

func (c *CreateInboundCommand) Execute(ctx context.Context, req *dto.CreateInboundRequest) (*dto.InboundResponse, error) {
	// Verify instance exists
	instance, err := c.instanceRepo.FindByID(ctx, req.XrayInstanceID)
	if err != nil {
		return nil, err
	}

	// Create inbound entity
	inbound, err := xray.NewInbound(
		req.XrayInstanceID,
		xray.InboundProtocol(req.Protocol),
		req.Port,
		xray.TransportType(req.Transport),
		xray.SecurityType(req.Security),
	)
	if err != nil {
		return nil, err
	}

	// Save to repository
	if err := c.inboundRepo.Create(ctx, inbound); err != nil {
		c.logger.Error("Failed to create inbound", "error", err)
		return nil, err
	}

	// Increment config version to trigger reload
	instance.IncrementConfigVersion()
	if err := c.instanceRepo.Update(ctx, instance); err != nil {
		c.logger.Error("Failed to update instance config version", "error", err, "instance_id", instance.ID)
		return nil, fmt.Errorf("failed to update instance: %w", err)
	}

	c.logger.Info("Inbound created successfully", "inbound_id", inbound.ID, "instance_id", req.XrayInstanceID)

	return mapper.ToInboundResponse(inbound), nil
}
