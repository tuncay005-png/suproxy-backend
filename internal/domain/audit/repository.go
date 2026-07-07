package audit

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// AuditFilters defines filter parameters for audit log queries
type AuditFilters struct {
	UserID     *uuid.UUID
	Action     *string
	EntityType *string
	EntityID   *uuid.UUID
	IPAddress  *string
	DateFrom   *time.Time
	DateTo     *time.Time
	Offset     int
	Limit      int
	SortBy     string
	SortOrder  string
}

type Repository interface {
	Create(ctx context.Context, log *Log) error
	FindByID(ctx context.Context, id uuid.UUID) (*Log, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]*Log, error)
	FindByEntityID(ctx context.Context, entityType string, entityID uuid.UUID) ([]*Log, error)
	List(ctx context.Context, offset, limit int) ([]*Log, error)
	ListWithFilters(ctx context.Context, filters AuditFilters) ([]*Log, int64, error)
	Count(ctx context.Context) (int64, error)
	CountByAction(ctx context.Context) (map[string]int64, error)
	CountByEntityType(ctx context.Context) (map[string]int64, error)
	CountUniqueUsers(ctx context.Context) (int64, error)
	CountUniqueIPs(ctx context.Context) (int64, error)
	GetOldestLogDate(ctx context.Context) (*time.Time, error)
	GetNewestLogDate(ctx context.Context) (*time.Time, error)
	DeleteOlderThan(ctx context.Context, date time.Time) error
}
