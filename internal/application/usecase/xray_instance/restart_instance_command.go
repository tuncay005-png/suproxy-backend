package xray_instance

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/xray"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type RestartInstanceCommand struct {
	instanceRepo xray.XrayInstanceRepository
	logger       *logger.Logger
}

func NewRestartInstanceCommand(
	instanceRepo xray.XrayInstanceRepository,
	logger *logger.Logger,
) *RestartInstanceCommand {
	return &RestartInstanceCommand{
		instanceRepo: instanceRepo,
		logger:       logger,
	}
}

func (c *RestartInstanceCommand) Execute(ctx context.Context, instanceID uuid.UUID) (*dto.XrayInstanceResponse, error) {
	// Find instance
	instance, err := c.instanceRepo.FindByID(ctx, instanceID)
	if err != nil {
		return nil, err
	}

	// Restart instance
	if err := instance.Restart(); err != nil {
		return nil, err
	}

	// Save changes
	if err := c.instanceRepo.Update(ctx, instance); err != nil {
		c.logger.Error("Failed to restart xray instance", "error", err, "instance_id", instanceID)
		return nil, err
	}

	c.logger.Info("Xray instance restarted successfully", "instance_id", instanceID)

	// TODO: In future, trigger actual Xray process restart via XrayManager
	// This will reload configuration with new ConfigVersion

	return mapper.ToXrayInstanceResponse(instance), nil
}
