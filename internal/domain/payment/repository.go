package payment

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, payment *Payment) error
	FindByID(ctx context.Context, id uuid.UUID) (*Payment, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]*Payment, error)
	FindBySubscriptionID(ctx context.Context, subscriptionID uuid.UUID) ([]*Payment, error)
	FindByTransactionID(ctx context.Context, transactionID string) (*Payment, error)
	Update(ctx context.Context, payment *Payment) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*Payment, error)
	Count(ctx context.Context) (int64, error)
}
