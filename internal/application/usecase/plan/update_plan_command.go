package plan

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/subscription"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type UpdatePlanCommand struct {
	planRepo subscription.PlanRepository
	logger   *logger.Logger
}

func NewUpdatePlanCommand(planRepo subscription.PlanRepository, logger *logger.Logger) *UpdatePlanCommand {
	return &UpdatePlanCommand{
		planRepo: planRepo,
		logger:   logger,
	}
}

func (c *UpdatePlanCommand) Execute(ctx context.Context, planID uuid.UUID, req *dto.UpdatePlanRequest) (*dto.PlanResponse, error) {
	// Find plan
	plan, err := c.planRepo.FindByID(ctx, planID)
	if err != nil {
		return nil, err
	}

	// Update details
	if req.Description != "" || req.TrafficLimitGB >= 0 || req.DeviceLimit > 0 || req.MaxSessions > 0 {
		description := req.Description
		if description == "" {
			description = plan.Description
		}

		trafficLimit := req.TrafficLimitGB
		if trafficLimit < 0 {
			trafficLimit = plan.TrafficLimitGB
		}

		deviceLimit := req.DeviceLimit
		if deviceLimit == 0 {
			deviceLimit = plan.DeviceLimit
		}

		maxSessions := req.MaxSessions
		if maxSessions == 0 {
			maxSessions = plan.MaxSessions
		}

		if err := plan.UpdateDetails(description, trafficLimit, deviceLimit, maxSessions); err != nil {
			return nil, err
		}
	}

	// Update price if provided
	if req.Price > 0 && req.Currency != "" {
		money, err := subscription.NewMoney(req.Price, req.Currency)
		if err != nil {
			return nil, err
		}
		plan.UpdatePrice(money, req.Currency)
	}

	// Update features if provided
	if req.Features != nil {
		// Clear existing features and add new ones
		plan.Features = []string{}
		for _, feature := range req.Features {
			plan.AddFeature(feature)
		}
	}

	// Save
	if err := c.planRepo.Update(ctx, plan); err != nil {
		c.logger.Error("Failed to update plan", "error", err, "plan_id", planID)
		return nil, err
	}

	c.logger.Info("Plan updated successfully", "plan_id", planID)

	return mapper.ToPlanResponse(plan), nil
}
