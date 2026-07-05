package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateServerRequest struct {
	Name        string   `json:"name" binding:"required"`
	Country     string   `json:"country" binding:"required"`
	City        string   `json:"city" binding:"required"`
	CountryCode string   `json:"country_code" binding:"required"`
	Latitude    float64  `json:"latitude"`
	Longitude   float64  `json:"longitude"`
	IPAddress   string   `json:"ip_address" binding:"required"`
	Domain      string   `json:"domain"`
	Port        int      `json:"port" binding:"required"`
	MaxUsers    int      `json:"max_users" binding:"required"`
	MaxBandwidth int64   `json:"max_bandwidth" binding:"required"`
	Tags        []string `json:"tags"`
}

type UpdateServerRequest struct {
	Name         string   `json:"name"`
	MaxUsers     int      `json:"max_users"`
	MaxBandwidth int64    `json:"max_bandwidth"`
	Tags         []string `json:"tags"`
}

type ServerResponse struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Country          string    `json:"country"`
	City             string    `json:"city"`
	CountryCode      string    `json:"country_code"`
	IPAddress        string    `json:"ip_address"`
	Domain           string    `json:"domain"`
	Port             int       `json:"port"`
	MaxUsers         int       `json:"max_users"`
	CurrentUsers     int       `json:"current_users"`
	MaxBandwidth     int64     `json:"max_bandwidth"`
	UsedBandwidth    int64     `json:"used_bandwidth"`
	Status           string    `json:"status"`
	UsagePercentage  float64   `json:"usage_percentage"`
	Tags             []string  `json:"tags"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type ServerListResponse struct {
	Servers []*ServerResponse `json:"servers"`
	Total   int64             `json:"total"`
	Offset  int               `json:"offset"`
	Limit   int               `json:"limit"`
}
