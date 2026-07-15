package dto

import (
	"time"

	"github.com/google/uuid"
)

// Admin Xray Instance List Request with filters
type AdminXrayInstanceListRequest struct {
	Offset    int    `form:"offset"`
	Limit     int    `form:"limit" binding:"max=100"`
	NodeID    string `form:"node_id"` // UUID as string
	Status    string `form:"status" binding:"omitempty,oneof=running stopped failed"`
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=created_at status"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// Admin Xray Instance Response
type AdminXrayInstanceResponse struct {
	ID            uuid.UUID  `json:"id"`
	NodeID        uuid.UUID  `json:"node_id"`
	Version       string     `json:"version"`
	Status        string     `json:"status"`
	ConfigVersion int        `json:"config_version"`
	StartedAt     *time.Time `json:"started_at"`
	StoppedAt     *time.Time `json:"stopped_at"`
	Uptime        int64      `json:"uptime_seconds"` // Uptime in seconds
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// Admin Xray Instance List Response
type AdminXrayInstanceListResponse struct {
	Instances []*AdminXrayInstanceResponse `json:"instances"`
	Total     int64                        `json:"total"`
	Offset    int                          `json:"offset"`
	Limit     int                          `json:"limit"`
}

// Xray Instance Control Response
type XrayInstanceControlResponse struct {
	Success  bool                       `json:"success"`
	Message  string                     `json:"message"`
	Instance *AdminXrayInstanceResponse `json:"instance"`
}

// Health Check Response
type HealthCheckResponse struct {
	Healthy   bool      `json:"healthy"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	CheckedAt time.Time `json:"checked_at"`
}

// Config Regenerate Response
type ConfigRegenerateResponse struct {
	Success       bool      `json:"success"`
	Message       string    `json:"message"`
	ConfigSize    int       `json:"config_size_bytes"`
	RegeneratedAt time.Time `json:"regenerated_at"`
	DurationMs    int64     `json:"duration_ms"`
}

// Instance Statistics Response
type InstanceStatsResponse struct {
	InstanceID      uuid.UUID `json:"instance_id"`
	TotalInbounds   int       `json:"total_inbounds"`
	EnabledInbounds int       `json:"enabled_inbounds"`
	TotalClients    int       `json:"total_clients"`
	EnabledClients  int       `json:"enabled_clients"`
	ConfigVersion   int       `json:"config_version"`
	Uptime          int64     `json:"uptime_seconds"`
}
