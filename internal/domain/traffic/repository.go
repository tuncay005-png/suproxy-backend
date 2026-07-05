package traffic

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, usage *Usage) error
	FindByID(ctx context.Context, id uuid.UUID) (*Usage, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]*Usage, error)
	FindBySubscriptionID(ctx context.Context, subscriptionID uuid.UUID, from, to time.Time) ([]*Usage, error)
	SumByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time) (int64, error)
	SumBySubscriptionID(ctx context.Context, subscriptionID uuid.UUID, from, to time.Time) (int64, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteOlderThan(ctx context.Context, date time.Time) error
}
