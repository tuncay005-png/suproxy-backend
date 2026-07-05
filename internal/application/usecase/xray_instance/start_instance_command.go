package xray_instance

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/xray"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type StartInstanceCommand struct {
	instanceRepo xray.XrayInstanceRepository
	logger       *logger.Logger
}

func NewStartInstanceCommand(
	instanceRepo xray.XrayInstanceRepository,
	logger *logger.Logger,
) *StartInstanceCommand {
	return &StartInstanceCommand{
		instanceRepo: instanceRepo,
		logger:       logger,
	}
}

func (c *StartInstanceCommand) Execute(ctx context.Context, instanceID uuid.UUID) (*dto.XrayInstanceResponse, error) {
	// Find instance
	instance, err := c.instanceRepo.FindByID(ctx, instanceID)
	if err != nil {
		return nil, err
	}

	// Start instance
	if err := instance.Start(); err != nil {
		return nil, err
	}

	// Save changes
	if err := c.instanceRepo.Update(ctx, instance); err != nil {
		c.logger.Error("Failed to start xray instance", "error", err, "instance_id", instanceID)
		return nil, err
	}

	c.logger.Info("Xray instance started successfully", "instance_id", instanceID)

	// TODO: In future, trigger actual Xray process start via XrayManager
	// XrayManager will use this domain state to manage the real process

	return mapper.ToXrayInstanceResponse(instance), nil
}
