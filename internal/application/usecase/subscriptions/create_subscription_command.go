package subscriptions

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/subscription"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type CreateSubscriptionCommand struct {
	subRepo  subscription.SubscriptionRepository
	planRepo subscription.PlanRepository
	logger   *logger.Logger
}

func NewCreateSubscriptionCommand(
	subRepo subscription.SubscriptionRepository,
	planRepo subscription.PlanRepository,
	logger *logger.Logger,
) *CreateSubscriptionCommand {
	return &CreateSubscriptionCommand{
		subRepo:  subRepo,
		planRepo: planRepo,
		logger:   logger,
	}
}

func (c *CreateSubscriptionCommand) Execute(ctx context.Context, userID uuid.UUID, req *dto.CreateSubscriptionRequest) (*dto.SubscriptionResponse, error) {
	// Check if user already has an active subscription
	existing, err := c.subRepo.FindActiveByUserID(ctx, userID)
	if err == nil && existing != nil {
		return nil, subscription.ErrUserAlreadyHasSubscription
	}

	// Find plan
	plan, err := c.planRepo.FindByID(ctx, req.PlanID)
	if err != nil {
		return nil, err
	}

	// Check if plan is active
	if !plan.IsActiveStatus() {
		return nil, subscription.ErrPlanNotActive
	}

	// Create subscription
	sub, err := subscription.NewSubscription(userID, plan.ID, plan, req.AutoRenew)
	if err != nil {
		return nil, err
	}

	// Save
	if err := c.subRepo.Create(ctx, sub); err != nil {
		c.logger.Error("Failed to create subscription", "error", err)
		return nil, err
	}

	c.logger.Info("Subscription created successfully", "subscription_id", sub.ID, "user_id", userID)

	return mapper.ToSubscriptionResponse(sub, plan), nil
}
