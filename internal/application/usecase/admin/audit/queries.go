package audit

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/audit"
)

// ListAuditLogsQuery handles listing audit logs with advanced filters
type ListAuditLogsQuery struct {
	auditRepo audit.Repository
}

func NewListAuditLogsQuery(auditRepo audit.Repository) *ListAuditLogsQuery {
	return &ListAuditLogsQuery{
		auditRepo: auditRepo,
	}
}

// ListAuditLogsParams defines parameters for listing audit logs
type ListAuditLogsParams struct {
	Offset     int
	Limit      int
	UserID     string
	Action     string
	EntityType string
	EntityID   string
	IPAddress  string
	DateFrom   string
	DateTo     string
	SortBy     string
	SortOrder  string
}

func (q *ListAuditLogsQuery) Execute(ctx context.Context, params ListAuditLogsParams) ([]*audit.Log, int64, error) {
	filters := audit.AuditFilters{
		Offset:    params.Offset,
		Limit:     params.Limit,
		SortBy:    params.SortBy,
		SortOrder: params.SortOrder,
	}

	// Parse user ID filter
	if params.UserID != "" {
		userID, err := uuid.Parse(params.UserID)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid user_id: %w", err)
		}
		filters.UserID = &userID
	}

	// Parse action filter
	if params.Action != "" {
		filters.Action = &params.Action
	}

	// Parse entity type filter
	if params.EntityType != "" {
		filters.EntityType = &params.EntityType
	}

	// Parse entity ID filter
	if params.EntityID != "" {
		entityID, err := uuid.Parse(params.EntityID)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid entity_id: %w", err)
		}
		filters.EntityID = &entityID
	}

	// Parse IP address filter
	if params.IPAddress != "" {
		filters.IPAddress = &params.IPAddress
	}

	// Parse date from filter
	if params.DateFrom != "" {
		dateFrom, err := time.Parse(time.RFC3339, params.DateFrom)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid date_from format (use RFC3339): %w", err)
		}
		filters.DateFrom = &dateFrom
	}

	// Parse date to filter
	if params.DateTo != "" {
		dateTo, err := time.Parse(time.RFC3339, params.DateTo)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid date_to format (use RFC3339): %w", err)
		}
		filters.DateTo = &dateTo
	}

	return q.auditRepo.ListWithFilters(ctx, filters)
}

// GetAuditLogQuery handles retrieving a single audit log by ID
type GetAuditLogQuery struct {
	auditRepo audit.Repository
}

func NewGetAuditLogQuery(auditRepo audit.Repository) *GetAuditLogQuery {
	return &GetAuditLogQuery{
		auditRepo: auditRepo,
	}
}

func (q *GetAuditLogQuery) Execute(ctx context.Context, logID uuid.UUID) (*audit.Log, error) {
	return q.auditRepo.FindByID(ctx, logID)
}

// GetAuditStatsQuery handles retrieving audit statistics
type GetAuditStatsQuery struct {
	auditRepo audit.Repository
}

func NewGetAuditStatsQuery(auditRepo audit.Repository) *GetAuditStatsQuery {
	return &GetAuditStatsQuery{
		auditRepo: auditRepo,
	}
}

type AuditStats struct {
	TotalLogs         int64
	LogsByAction      map[string]int64
	LogsByEntityType  map[string]int64
	UniqueUsers       int64
	UniqueIPAddresses int64
	OldestLog         *time.Time
	NewestLog         *time.Time
}

func (q *GetAuditStatsQuery) Execute(ctx context.Context) (*AuditStats, error) {
	// Get total count
	total, err := q.auditRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count logs: %w", err)
	}

	// Get counts by action
	logsByAction, err := q.auditRepo.CountByAction(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count by action: %w", err)
	}

	// Get counts by entity type
	logsByEntityType, err := q.auditRepo.CountByEntityType(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count by entity type: %w", err)
	}

	// Get unique users count
	uniqueUsers, err := q.auditRepo.CountUniqueUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count unique users: %w", err)
	}

	// Get unique IPs count
	uniqueIPs, err := q.auditRepo.CountUniqueIPs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count unique IPs: %w", err)
	}

	// Get oldest log date
	oldestLog, err := q.auditRepo.GetOldestLogDate(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get oldest log date: %w", err)
	}

	// Get newest log date
	newestLog, err := q.auditRepo.GetNewestLogDate(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get newest log date: %w", err)
	}

	return &AuditStats{
		TotalLogs:         total,
		LogsByAction:      logsByAction,
		LogsByEntityType:  logsByEntityType,
		UniqueUsers:       uniqueUsers,
		UniqueIPAddresses: uniqueIPs,
		OldestLog:         oldestLog,
		NewestLog:         newestLog,
	}, nil
}
