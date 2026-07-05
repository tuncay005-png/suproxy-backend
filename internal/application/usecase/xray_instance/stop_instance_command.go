package xray_instance

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/xray"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type StopInstanceCommand struct {
	instanceRepo xray.XrayInstanceRepository
	logger       *logger.Logger
}

func NewStopInstanceCommand(
	instanceRepo xray.XrayInstanceRepository,
	logger *logger.Logger,
) *StopInstanceCommand {
	return &StopInstanceCommand{
		instanceRepo: instanceRepo,
		logger:       logger,
	}
}

func (c *StopInstanceCommand) Execute(ctx context.Context, instanceID uuid.UUID) (*dto.XrayInstanceResponse, error) {
	// Find instance
	instance, err := c.instanceRepo.FindByID(ctx, instanceID)
	if err != nil {
		return nil, err
	}

	// Stop instance
	if err := instance.Stop(); err != nil {
		return nil, err
	}

	// Save changes
	if err := c.instanceRepo.Update(ctx, instance); err != nil {
		c.logger.Error("Failed to stop xray instance", "error", err, "instance_id", instanceID)
		return nil, err
	}

	c.logger.Info("Xray instance stopped successfully", "instance_id", instanceID)

	// TODO: In future, trigger actual Xray process stop via XrayManager

	return mapper.ToXrayInstanceResponse(instance), nil
}
