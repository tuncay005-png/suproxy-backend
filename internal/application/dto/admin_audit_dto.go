package dto

import (
	"time"

	"github.com/google/uuid"
)

// AdminAuditListRequest defines request parameters for listing audit logs
type AdminAuditListRequest struct {
	// Pagination
	Offset int `form:"offset"`
	Limit  int `form:"limit"`

	// Filters
	UserID     string `form:"user_id"`     // Filter by user ID
	Action     string `form:"action"`      // Filter by action type (create, update, delete, etc.)
	EntityType string `form:"entity_type"` // Filter by entity type (user, xray_instance, etc.)
	EntityID   string `form:"entity_id"`   // Filter by specific entity ID
	IPAddress  string `form:"ip_address"`  // Filter by IP address
	DateFrom   string `form:"date_from"`   // Filter by start date (RFC3339 format)
	DateTo     string `form:"date_to"`     // Filter by end date (RFC3339 format)

	// Sorting
	SortBy    string `form:"sort_by"`    // Sort field (created_at, action, entity_type)
	SortOrder string `form:"sort_order"` // Sort order (asc, desc)
}

// AdminAuditResponse represents a single audit log entry
type AdminAuditResponse struct {
	ID         uuid.UUID              `json:"id"`
	UserID     uuid.UUID              `json:"user_id"`
	Action     string                 `json:"action"`
	EntityType string                 `json:"entity_type"`
	EntityID   uuid.UUID              `json:"entity_id"`
	IPAddress  string                 `json:"ip_address"`
	UserAgent  string                 `json:"user_agent"`
	Metadata   map[string]interface{} `json:"metadata"`
	CreatedAt  time.Time              `json:"created_at"`
}

// AdminAuditListResponse represents paginated audit log list
type AdminAuditListResponse struct {
	Data       []AdminAuditResponse `json:"data"`
	Total      int64                `json:"total"`
	Offset     int                  `json:"offset"`
	Limit      int                  `json:"limit"`
	TotalPages int                  `json:"total_pages"`
}

// AuditStatsResponse represents audit statistics
type AuditStatsResponse struct {
	TotalLogs         int64            `json:"total_logs"`
	LogsByAction      map[string]int64 `json:"logs_by_action"`
	LogsByEntityType  map[string]int64 `json:"logs_by_entity_type"`
	UniqueUsers       int64            `json:"unique_users"`
	UniqueIPAddresses int64            `json:"unique_ip_addresses"`
	OldestLog         *time.Time       `json:"oldest_log,omitempty"`
	NewestLog         *time.Time       `json:"newest_log,omitempty"`
}
