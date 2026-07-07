package dto

import (
	"time"

	"github.com/google/uuid"
)

// Admin Inbound List Request with filters
type AdminInboundListRequest struct {
	Offset     int    `form:"offset"`
	Limit      int    `form:"limit" binding:"max=100"`
	InstanceID string `form:"instance_id"` // UUID as string
	Protocol   string `form:"protocol" binding:"omitempty,oneof=vless vmess trojan shadowsocks"`
	Enabled    string `form:"enabled" binding:"omitempty,oneof=true false"`
	SortBy     string `form:"sort_by" binding:"omitempty,oneof=created_at port"`
	SortOrder  string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// Admin Inbound Response
type AdminInboundResponse struct {
	ID             uuid.UUID `json:"id"`
	XrayInstanceID uuid.UUID `json:"xray_instance_id"`
	Protocol       string    `json:"protocol"`
	Port           int       `json:"port"`
	Transport      string    `json:"transport"`
	Security       string    `json:"security"`
	Enabled        bool      `json:"enabled"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Admin Inbound List Response
type AdminInboundListResponse struct {
	Inbounds []*AdminInboundResponse `json:"inbounds"`
	Total    int64                   `json:"total"`
	Offset   int                     `json:"offset"`
	Limit    int                     `json:"limit"`
}

// Create Inbound Request (Admin)
type AdminCreateInboundRequest struct {
	XrayInstanceID string `json:"xray_instance_id" binding:"required,uuid"`
	Protocol       string `json:"protocol" binding:"required,oneof=vless vmess trojan shadowsocks"`
	Port           int    `json:"port" binding:"required,min=1,max=65535"`
	Transport      string `json:"transport" binding:"required,oneof=tcp ws grpc"`
	Security       string `json:"security" binding:"required,oneof=none tls reality"`
}

// Update Inbound Request (Admin)
type AdminUpdateInboundRequest struct {
	Port      int    `json:"port" binding:"omitempty,min=1,max=65535"`
	Transport string `json:"transport" binding:"omitempty,oneof=tcp ws grpc"`
	Security  string `json:"security" binding:"omitempty,oneof=none tls reality"`
}

// Inbound Operation Response
type InboundOperationResponse struct {
	Success  bool                  `json:"success"`
	Message  string                `json:"message"`
	Inbound  *AdminInboundResponse `json:"inbound"`
}
