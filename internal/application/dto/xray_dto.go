package dto

import (
	"time"

	"github.com/google/uuid"
)

// XrayInstance DTOs
type CreateXrayInstanceRequest struct {
	NodeID  uuid.UUID `json:"node_id" binding:"required"`
	Version string    `json:"version" binding:"required"`
}

type UpdateXrayInstanceRequest struct {
	Version string `json:"version"`
}

type XrayInstanceResponse struct {
	ID            uuid.UUID  `json:"id"`
	NodeID        uuid.UUID  `json:"node_id"`
	Version       string     `json:"version"`
	Status        string     `json:"status"`
	ConfigVersion int        `json:"config_version"`
	StartedAt     *time.Time `json:"started_at"`
	StoppedAt     *time.Time `json:"stopped_at"`
	Uptime        int64      `json:"uptime_seconds"`
	IsRunning     bool       `json:"is_running"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type XrayInstanceListResponse struct {
	Instances []*XrayInstanceResponse `json:"instances"`
	Total     int64                   `json:"total"`
	Offset    int                     `json:"offset"`
	Limit     int                     `json:"limit"`
}

// Inbound DTOs
type CreateInboundRequest struct {
	XrayInstanceID uuid.UUID `json:"xray_instance_id" binding:"required"`
	Protocol       string    `json:"protocol" binding:"required"`
	Port           int       `json:"port" binding:"required"`
	Transport      string    `json:"transport" binding:"required"`
	Security       string    `json:"security" binding:"required"`
}

type UpdateInboundRequest struct {
	Port      *int   `json:"port"`
	Transport string `json:"transport"`
	Security  string `json:"security"`
}

type InboundResponse struct {
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

type InboundListResponse struct {
	Inbounds []*InboundResponse `json:"inbounds"`
	Total    int64              `json:"total"`
	Offset   int                `json:"offset"`
	Limit    int                `json:"limit"`
}

// Client DTOs
type CreateClientRequest struct {
	InboundID uuid.UUID `json:"inbound_id" binding:"required"`
	UserID    uuid.UUID `json:"user_id" binding:"required"`
	UUID      string    `json:"uuid" binding:"required"`
	Flow      string    `json:"flow"`
	Email     string    `json:"email" binding:"required"`
}

type UpdateClientRequest struct {
	Flow  string `json:"flow"`
	Email string `json:"email"`
}

type ClientResponse struct {
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

type ClientListResponse struct {
	Clients []*ClientResponse `json:"clients"`
	Total   int64             `json:"total"`
	Offset  int               `json:"offset"`
	Limit   int               `json:"limit"`
}

// Reality DTOs
type CreateRealityConfigRequest struct {
	InboundID   uuid.UUID `json:"inbound_id" binding:"required"`
	PrivateKey  string    `json:"private_key" binding:"required"`
	PublicKey   string    `json:"public_key" binding:"required"`
	ShortID     string    `json:"short_id"`
	ServerName  string    `json:"server_name" binding:"required"`
	Fingerprint string    `json:"fingerprint" binding:"required"`
	SpiderX     string    `json:"spider_x"`
}

type UpdateRealityConfigRequest struct {
	ServerName  string `json:"server_name"`
	Fingerprint string `json:"fingerprint"`
	ShortID     string `json:"short_id"`
	SpiderX     string `json:"spider_x"`
}

type RegenerateRealityKeysRequest struct {
	PrivateKey string `json:"private_key" binding:"required"`
	PublicKey  string `json:"public_key" binding:"required"`
}

type RealityConfigResponse struct {
	ID          uuid.UUID `json:"id"`
	InboundID   uuid.UUID `json:"inbound_id"`
	PrivateKey  string    `json:"private_key"`
	PublicKey   string    `json:"public_key"`
	ShortID     string    `json:"short_id"`
	ServerName  string    `json:"server_name"`
	Fingerprint string    `json:"fingerprint"`
	SpiderX     string    `json:"spider_x"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
