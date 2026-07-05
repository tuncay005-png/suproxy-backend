package dto

import (
	"time"

	"github.com/google/uuid"
)

type SessionInfo struct {
	ID         uuid.UUID  `json:"id"`
	DeviceName string     `json:"device_name"`
	Platform   string     `json:"platform"`
	IPAddress  string     `json:"ip_address"`
	LastUsedAt *time.Time `json:"last_used_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

type ActiveSessionsResponse struct {
	Sessions []*SessionInfo `json:"sessions"`
	Total    int            `json:"total"`
}
