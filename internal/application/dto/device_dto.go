package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateDeviceRequest struct {
	Name       string `json:"name" binding:"required"`
	DeviceType string `json:"device_type" binding:"required"`
	Identifier string `json:"identifier" binding:"required"`
}

type UpdateDeviceRequest struct {
	Name string `json:"name" binding:"required"`
}

type DeviceResponse struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	Name          string     `json:"name"`
	DeviceType    string     `json:"device_type"`
	Identifier    string     `json:"identifier"`
	Status        string     `json:"status"`
	LastSeenAt    *time.Time `json:"last_seen_at"`
	LastIPAddress string     `json:"last_ip_address,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type DeviceListResponse struct {
	Devices []*DeviceResponse `json:"devices"`
	Total   int64             `json:"total"`
}
