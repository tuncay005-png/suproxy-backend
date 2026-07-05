package subscriptions

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/subscription"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type ExtendSubscriptionCommand struct {
	subRepo  subscription.SubscriptionRepository
	planRepo subscription.PlanRepository
	logger   *logger.Logger
}

func NewExtendSubscriptionCommand(
	subRepo subscription.SubscriptionRepository,
	planRepo subscription.PlanRepository,
	logger *logger.Logger,
) *ExtendSubscriptionCommand {
	return &ExtendSubscriptionCommand{
		subRepo:  subRepo,
		planRepo: planRepo,
		logger:   logger,
	}
}

func (c *ExtendSubscriptionCommand) Execute(ctx context.Context, subscriptionID uuid.UUID, req *dto.ExtendSubscriptionRequest) (*dto.SubscriptionResponse, error) {
	// Find subscription
	sub, err := c.subRepo.FindByID(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Extend subscription
	if err := sub.Extend(req.Days); err != nil {
		return nil, err
	}

	// Save changes
	if err := c.subRepo.Update(ctx, sub); err != nil {
		c.logger.Error("Failed to extend subscription", "error", err, "subscription_id", subscriptionID)
		return nil, err
	}

	// Load plan for response
	plan, err := c.planRepo.FindByID(ctx, sub.PlanID)
	if err != nil {
		return nil, err
	}

	c.logger.Info("Subscription extended successfully", "subscription_id", subscriptionID, "days", req.Days)

	return mapper.ToSubscriptionResponse(sub, plan), nil
}
