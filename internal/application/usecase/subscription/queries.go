package subscription

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
)

type GetPlanQuery interface {
	Execute(ctx context.Context, planID uuid.UUID) (*dto.PlanResponse, error)
}

type ListPlansQuery interface {
	Execute(ctx context.Context) ([]*dto.PlanResponse, error)
}

type ListActivePlansQuery interface {
	Execute(ctx context.Context) ([]*dto.PlanResponse, error)
}

type GetSubscriptionQuery interface {
	Execute(ctx context.Context, subscriptionID uuid.UUID) (*dto.SubscriptionResponse, error)
}

type GetUserSubscriptionQuery interface {
	Execute(ctx context.Context, userID uuid.UUID) (*dto.SubscriptionResponse, error)
}
