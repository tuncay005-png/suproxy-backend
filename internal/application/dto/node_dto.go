package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateNodeRequest struct {
	ServerID         uuid.UUID `json:"server_id" binding:"required"`
	Protocol         string    `json:"protocol" binding:"required"`
	Port             int       `json:"port" binding:"required"`
	MaxUsers         int       `json:"max_users" binding:"required"`
	BandwidthLimitGB int64     `json:"bandwidth_limit_gb" binding:"required"`
}

type UpdateNodeRequest struct {
	MaxUsers         *int   `json:"max_users"`
	BandwidthLimitGB *int64 `json:"bandwidth_limit_gb"`
	Version          string `json:"version"`
}

type UpdateNodeMetricsRequest struct {
	CPUUsage  float64 `json:"cpu_usage" binding:"required"`
	RAMUsage  float64 `json:"ram_usage" binding:"required"`
	LatencyMs int     `json:"latency_ms" binding:"required"`
}

type NodeResponse struct {
	ID                    uuid.UUID `json:"id"`
	ServerID              uuid.UUID `json:"server_id"`
	Protocol              string    `json:"protocol"`
	Port                  int       `json:"port"`
	MaxUsers              int       `json:"max_users"`
	CurrentUsers          int       `json:"current_users"`
	AvailableSlots        int       `json:"available_slots"`
	UserLoadPercentage    float64   `json:"user_load_percentage"`
	BandwidthLimitBytes   int64     `json:"bandwidth_limit_bytes"`
	BandwidthUsedBytes    int64     `json:"bandwidth_used_bytes"`
	BandwidthLimitGB      float64   `json:"bandwidth_limit_gb"`
	BandwidthUsedGB       float64   `json:"bandwidth_used_gb"`
	RemainingBandwidthGB  float64   `json:"remaining_bandwidth_gb"`
	BandwidthUsagePercent float64   `json:"bandwidth_usage_percent"`
	HasUnlimitedBandwidth bool      `json:"has_unlimited_bandwidth"`
	CPUUsage              float64   `json:"cpu_usage"`
	RAMUsage              float64   `json:"ram_usage"`
	LatencyMs             int       `json:"latency_ms"`
	Version               string    `json:"version"`
	HealthStatus          string    `json:"health_status"`
	IsHealthy             bool      `json:"is_healthy"`
	IsOverloaded          bool      `json:"is_overloaded"`
	CanAcceptUser         bool      `json:"can_accept_user"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type NodeListResponse struct {
	Nodes  []*NodeResponse `json:"nodes"`
	Total  int64           `json:"total"`
	Offset int             `json:"offset"`
	Limit  int             `json:"limit"`
}
