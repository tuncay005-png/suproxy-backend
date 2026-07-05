package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreatePlanRequest struct {
	Name           string   `json:"name" binding:"required"`
	Description    string   `json:"description"`
	DurationDays   int      `json:"duration_days" binding:"required"`
	TrafficLimitGB int64    `json:"traffic_limit_gb" binding:"required"`
	DeviceLimit    int      `json:"device_limit" binding:"required"`
	MaxSessions    int      `json:"max_sessions" binding:"required"`
	Price          int64    `json:"price" binding:"required"`
	Currency       string   `json:"currency" binding:"required"`
	Features       []string `json:"features"`
}

type UpdatePlanRequest struct {
	Description    string   `json:"description"`
	TrafficLimitGB int64    `json:"traffic_limit_gb"`
	DeviceLimit    int      `json:"device_limit"`
	MaxSessions    int      `json:"max_sessions"`
	Price          int64    `json:"price"`
	Currency       string   `json:"currency"`
	Features       []string `json:"features"`
}

type PlanResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	DurationDays   int       `json:"duration_days"`
	TrafficLimitGB int64     `json:"traffic_limit_gb"`
	DeviceLimit    int       `json:"device_limit"`
	MaxSessions    int       `json:"max_sessions"`
	Price          int64     `json:"price"`
	Currency       string    `json:"currency"`
	IsActive       bool      `json:"is_active"`
	Features       []string  `json:"features"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type PlanListResponse struct {
	Plans  []*PlanResponse `json:"plans"`
	Total  int64           `json:"total"`
}

type CreateSubscriptionRequest struct {
	PlanID    uuid.UUID `json:"plan_id" binding:"required"`
	AutoRenew bool      `json:"auto_renew"`
}

type ExtendSubscriptionRequest struct {
	Days int `json:"days" binding:"required,min=1"`
}

type SubscriptionResponse struct {
	ID                   uuid.UUID     `json:"id"`
	UserID               uuid.UUID     `json:"user_id"`
	Plan                 *PlanResponse `json:"plan"`
	Status               string        `json:"status"`
	StartedAt            time.Time     `json:"started_at"`
	ExpiresAt            time.Time     `json:"expires_at"`
	TrafficUsedBytes     int64         `json:"traffic_used_bytes"`
	TrafficLimitBytes    int64         `json:"traffic_limit_bytes"`
	TrafficUsedGB        float64       `json:"traffic_used_gb"`
	TrafficLimitGB       float64       `json:"traffic_limit_gb"`
	RemainingTrafficGB   float64       `json:"remaining_traffic_gb"`
	TrafficUsagePercent  float64       `json:"traffic_usage_percent"`
	DaysRemaining        int           `json:"days_remaining"`
	AutoRenew            bool          `json:"auto_renew"`
	CanConnect           bool          `json:"can_connect"`
	HasUnlimitedTraffic  bool          `json:"has_unlimited_traffic"`
	CreatedAt            time.Time     `json:"created_at"`
	UpdatedAt            time.Time     `json:"updated_at"`
}
