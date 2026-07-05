package subscription

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
)

type CreatePlanCommand interface {
	Execute(ctx context.Context, req *dto.CreatePlanRequest) (*dto.PlanResponse, error)
}

type UpdatePlanCommand interface {
	Execute(ctx context.Context, planID uuid.UUID, req *dto.UpdatePlanRequest) (*dto.PlanResponse, error)
}

type ActivatePlanCommand interface {
	Execute(ctx context.Context, planID uuid.UUID) error
}

type DeactivatePlanCommand interface {
	Execute(ctx context.Context, planID uuid.UUID) error
}

type CreateSubscriptionCommand interface {
	Execute(ctx context.Context, userID uuid.UUID, req *dto.CreateSubscriptionRequest) (*dto.SubscriptionResponse, error)
}

type CancelSubscriptionCommand interface {
	Execute(ctx context.Context, subscriptionID uuid.UUID) error
}

type RenewSubscriptionCommand interface {
	Execute(ctx context.Context, subscriptionID uuid.UUID, period string) error
}

type UpgradeSubscriptionCommand interface {
	Execute(ctx context.Context, subscriptionID, planID uuid.UUID) error
}
