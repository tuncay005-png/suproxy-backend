package subscriptions

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/subscription"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type GetSubscriptionQuery struct {
	subRepo  subscription.SubscriptionRepository
	planRepo subscription.PlanRepository
	logger   *logger.Logger
}

func NewGetSubscriptionQuery(
	subRepo subscription.SubscriptionRepository,
	planRepo subscription.PlanRepository,
	logger *logger.Logger,
) *GetSubscriptionQuery {
	return &GetSubscriptionQuery{
		subRepo:  subRepo,
		planRepo: planRepo,
		logger:   logger,
	}
}

func (q *GetSubscriptionQuery) ExecuteByID(ctx context.Context, subscriptionID uuid.UUID) (*dto.SubscriptionResponse, error) {
	sub, err := q.subRepo.FindByID(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	plan, err := q.planRepo.FindByID(ctx, sub.PlanID)
	if err != nil {
		return nil, err
	}

	return mapper.ToSubscriptionResponse(sub, plan), nil
}

func (q *GetSubscriptionQuery) ExecuteByUserID(ctx context.Context, userID uuid.UUID) (*dto.SubscriptionResponse, error) {
	sub, err := q.subRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	plan, err := q.planRepo.FindByID(ctx, sub.PlanID)
	if err != nil {
		return nil, err
	}

	return mapper.ToSubscriptionResponse(sub, plan), nil
}
