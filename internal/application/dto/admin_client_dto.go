package dto

import (
	"time"

	"github.com/google/uuid"
)

// Admin Client List Request with filters
type AdminClientListRequest struct {
	Offset    int    `form:"offset"`
	Limit     int    `form:"limit" binding:"max=100"`
	InboundID string `form:"inbound_id"` // UUID as string
	UserID    string `form:"user_id"`    // UUID as string
	Enabled   string `form:"enabled" binding:"omitempty,oneof=true false"`
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=created_at email"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// Admin Client Response
type AdminClientResponse struct {
	ID        uuid.UUID `json:"id"`
	InboundID uuid.UUID `json:"inbound_id"`
	UserID    uuid.UUID `json:"user_id"`
	UUID      string    `json:"uuid"`
	Flow      string    `json:"flow"`
	Email     string    `json:"email"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Admin Client List Response
type AdminClientListResponse struct {
	Clients []*AdminClientResponse `json:"clients"`
	Total   int64                  `json:"total"`
	Offset  int                    `json:"offset"`
	Limit   int                    `json:"limit"`
}

// Create Client Request (Manual)
type AdminCreateClientRequest struct {
	InboundID string `json:"inbound_id" binding:"required,uuid"`
	UserID    string `json:"user_id" binding:"required,uuid"`
	Email     string `json:"email" binding:"required,email"`
	Flow      string `json:"flow" binding:"omitempty"`
}

// Regenerate UUID Request
type RegenerateUUIDRequest struct {
	// Empty - UUID will be auto-generated
}

// Client Operation Response
type ClientOperationResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Client  *AdminClientResponse `json:"client"`
}
