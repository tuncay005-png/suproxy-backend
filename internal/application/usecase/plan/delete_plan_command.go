package plan

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/subscription"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type DeletePlanCommand struct {
	planRepo subscription.PlanRepository
	logger   *logger.Logger
}

func NewDeletePlanCommand(planRepo subscription.PlanRepository, logger *logger.Logger) *DeletePlanCommand {
	return &DeletePlanCommand{
		planRepo: planRepo,
		logger:   logger,
	}
}

func (c *DeletePlanCommand) Execute(ctx context.Context, planID uuid.UUID) error {
	// Check if plan exists
	_, err := c.planRepo.FindByID(ctx, planID)
	if err != nil {
		return err
	}

	// Delete plan
	if err := c.planRepo.Delete(ctx, planID); err != nil {
		c.logger.Error("Failed to delete plan", "error", err, "plan_id", planID)
		return err
	}

	c.logger.Info("Plan deleted successfully", "plan_id", planID)
	return nil
}
