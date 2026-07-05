package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreatePlanRequest struct {
	Name         string   `json:"name" binding:"required"`
	Description  string   `json:"description"`
	Price        int64    `json:"price" binding:"required"`
	Currency     string   `json:"currency" binding:"required"`
	TrafficLimit int64    `json:"traffic_limit" binding:"required"`
	DeviceLimit  int      `json:"device_limit" binding:"required"`
	Features     []string `json:"features"`
}

type UpdatePlanRequest struct {
	Description  string   `json:"description"`
	Price        int64    `json:"price"`
	Currency     string   `json:"currency"`
	TrafficLimit int64    `json:"traffic_limit"`
	DeviceLimit  int      `json:"device_limit"`
	Features     []string `json:"features"`
}

type PlanResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Price        int64     `json:"price"`
	Currency     string    `json:"currency"`
	TrafficLimit int64     `json:"traffic_limit"`
	DeviceLimit  int       `json:"device_limit"`
	Features     []string  `json:"features"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateSubscriptionRequest struct {
	PlanID    uuid.UUID `json:"plan_id" binding:"required"`
	Period    string    `json:"period" binding:"required"`
	AutoRenew bool      `json:"auto_renew"`
}

type SubscriptionResponse struct {
	ID            uuid.UUID     `json:"id"`
	UserID        uuid.UUID     `json:"user_id"`
	Plan          *PlanResponse `json:"plan"`
	Status        string        `json:"status"`
	Period        string        `json:"period"`
	StartDate     time.Time     `json:"start_date"`
	EndDate       time.Time     `json:"end_date"`
	AutoRenew     bool          `json:"auto_renew"`
	DaysRemaining int           `json:"days_remaining"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}
