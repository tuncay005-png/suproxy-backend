package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateSessionRequest struct {
	DeviceID  uuid.UUID `json:"device_id" binding:"required"`
	NodeID    uuid.UUID `json:"node_id" binding:"required"`
	IPAddress string    `json:"ip_address"`
	Port      int       `json:"port"`
	Protocol  string    `json:"protocol"`
	UserAgent string    `json:"user_agent"`
}

type UpdateSessionTrafficRequest struct {
	BytesIn  int64 `json:"bytes_in"`
	BytesOut int64 `json:"bytes_out"`
}

type SessionResponse struct {
	ID             uuid.UUID  `json:"id"`
	UserID         uuid.UUID  `json:"user_id"`
	DeviceID       uuid.UUID  `json:"device_id"`
	NodeID         uuid.UUID  `json:"node_id"`
	SubscriptionID uuid.UUID  `json:"subscription_id"`
	Status         string     `json:"status"`
	ConnectedAt    time.Time  `json:"connected_at"`
	DisconnectedAt *time.Time `json:"disconnected_at"`
	Duration       int        `json:"duration"`
	BytesIn        int64      `json:"bytes_in"`
	BytesOut       int64      `json:"bytes_out"`
	TotalBytes     int64      `json:"total_bytes"`
	CreatedAt      time.Time  `json:"created_at"`
}

type SessionListResponse struct {
	Sessions []*SessionResponse `json:"sessions"`
	Total    int64              `json:"total"`
}
