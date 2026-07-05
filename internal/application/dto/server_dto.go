package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateServerRequest struct {
	Name     string `json:"name" binding:"required"`
	Country  string `json:"country" binding:"required"`
	City     string `json:"city" binding:"required"`
	Hostname string `json:"hostname" binding:"required"`
	Provider string `json:"provider" binding:"required"`
	IPv4     string `json:"ipv4" binding:"required"`
	IPv6     string `json:"ipv6"`
	IsPublic bool   `json:"is_public"`
}

type UpdateServerRequest struct {
	Name     string `json:"name"`
	City     string `json:"city"`
	Provider string `json:"provider"`
	IPv6     string `json:"ipv6"`
	IsPublic *bool  `json:"is_public"`
}

type ServerResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Country   string    `json:"country"`
	City      string    `json:"city"`
	Hostname  string    `json:"hostname"`
	Provider  string    `json:"provider"`
	IPv4      string    `json:"ipv4"`
	IPv6      string    `json:"ipv6"`
	Status    string    `json:"status"`
	IsPublic  bool      `json:"is_public"`
	IsOnline  bool      `json:"is_online"`
	NodeCount int       `json:"node_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ServerListResponse struct {
	Servers []*ServerResponse `json:"servers"`
	Total   int64             `json:"total"`
	Offset  int               `json:"offset"`
	Limit   int               `json:"limit"`
}
