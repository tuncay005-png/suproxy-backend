package plan

import (
	"context"

	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/subscription"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type CreatePlanCommand struct {
	planRepo subscription.PlanRepository
	logger   *logger.Logger
}

func NewCreatePlanCommand(planRepo subscription.PlanRepository, logger *logger.Logger) *CreatePlanCommand {
	return &CreatePlanCommand{
		planRepo: planRepo,
		logger:   logger,
	}
}

func (c *CreatePlanCommand) Execute(ctx context.Context, req *dto.CreatePlanRequest) (*dto.PlanResponse, error) {
	// Create money value object
	money, err := subscription.NewMoney(req.Price, req.Currency)
	if err != nil {
		return nil, err
	}

	// Create plan entity
	plan, err := subscription.NewPlan(
		req.Name,
		req.Description,
		req.DurationDays,
		req.TrafficLimitGB,
		req.DeviceLimit,
		req.MaxSessions,
		money,
		req.Currency,
	)
	if err != nil {
		return nil, err
	}

	// Add features
	for _, feature := range req.Features {
		plan.AddFeature(feature)
	}

	// Save to repository
	if err := c.planRepo.Create(ctx, plan); err != nil {
		c.logger.Error("Failed to create plan", "error", err)
		return nil, err
	}

	c.logger.Info("Plan created successfully", "plan_id", plan.ID, "name", plan.Name)

	return mapper.ToPlanResponse(plan), nil
}
