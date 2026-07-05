package audit

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, log *Log) error
	FindByID(ctx context.Context, id uuid.UUID) (*Log, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]*Log, error)
	FindByEntityID(ctx context.Context, entityType string, entityID uuid.UUID) ([]*Log, error)
	List(ctx context.Context, offset, limit int) ([]*Log, error)
	Count(ctx context.Context) (int64, error)
	DeleteOlderThan(ctx context.Context, date time.Time) error
}
