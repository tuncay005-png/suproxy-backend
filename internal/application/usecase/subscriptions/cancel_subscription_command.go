package subscriptions

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/subscription"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type CancelSubscriptionCommand struct {
	subRepo  subscription.SubscriptionRepository
	planRepo subscription.PlanRepository
	logger   *logger.Logger
}

func NewCancelSubscriptionCommand(
	subRepo subscription.SubscriptionRepository,
	planRepo subscription.PlanRepository,
	logger *logger.Logger,
) *CancelSubscriptionCommand {
	return &CancelSubscriptionCommand{
		subRepo:  subRepo,
		planRepo: planRepo,
		logger:   logger,
	}
}

func (c *CancelSubscriptionCommand) Execute(ctx context.Context, subscriptionID uuid.UUID) (*dto.SubscriptionResponse, error) {
	// Find subscription
	sub, err := c.subRepo.FindByID(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Cancel subscription
	if err := sub.Cancel(); err != nil {
		return nil, err
	}

	// Save changes
	if err := c.subRepo.Update(ctx, sub); err != nil {
		c.logger.Error("Failed to cancel subscription", "error", err, "subscription_id", subscriptionID)
		return nil, err
	}

	// Load plan for response
	plan, err := c.planRepo.FindByID(ctx, sub.PlanID)
	if err != nil {
		return nil, err
	}

	c.logger.Info("Subscription cancelled successfully", "subscription_id", subscriptionID)

	return mapper.ToSubscriptionResponse(sub, plan), nil
}
