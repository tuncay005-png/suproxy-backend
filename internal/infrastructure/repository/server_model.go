package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/server"
)

type ServerModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name      string    `gorm:"type:varchar(255);not null;index"`
	Country   string    `gorm:"type:varchar(100);not null;index"`
	City      string    `gorm:"type:varchar(100);not null"`
	Hostname  string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Provider  string    `gorm:"type:varchar(100);not null"`
	IPv4      string    `gorm:"type:varchar(45);not null"`
	IPv6      string    `gorm:"type:varchar(45)"`
	Status    string    `gorm:"type:varchar(20);not null;default:'offline';index"`
	IsPublic  bool      `gorm:"not null;default:true;index"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (ServerModel) TableName() string {
	return "servers"
}

func toServerModel(s *server.Server) *ServerModel {
	return &ServerModel{
		ID:        s.ID,
		Name:      s.Name,
		Country:   s.Country,
		City:      s.City,
		Hostname:  s.Hostname,
		Provider:  s.Provider,
		IPv4:      s.IPv4,
		IPv6:      s.IPv6,
		Status:    string(s.Status),
		IsPublic:  s.IsPublic,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func toDomainServer(m *ServerModel) *server.Server {
	return &server.Server{
		ID:        m.ID,
		Name:      m.Name,
		Country:   m.Country,
		City:      m.City,
		Hostname:  m.Hostname,
		Provider:  m.Provider,
		IPv4:      m.IPv4,
		IPv6:      m.IPv6,
		Status:    server.Status(m.Status),
		IsPublic:  m.IsPublic,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
