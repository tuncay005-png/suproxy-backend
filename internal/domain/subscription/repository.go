package subscription

import (
	"context"

	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *Subscription) error
	FindByID(ctx context.Context, id uuid.UUID) (*Subscription, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) (*Subscription, error)
	FindActiveByUserID(ctx context.Context, userID uuid.UUID) (*Subscription, error)
	Update(ctx context.Context, subscription *Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*Subscription, error)
	Count(ctx context.Context) (int64, error)
}

type PlanRepository interface {
	Create(ctx context.Context, plan *Plan) error
	FindByID(ctx context.Context, id uuid.UUID) (*Plan, error)
	FindByName(ctx context.Context, name string) (*Plan, error)
	FindActive(ctx context.Context) ([]*Plan, error)
	Update(ctx context.Context, plan *Plan) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*Plan, error)
	Count(ctx context.Context) (int64, error)
}
