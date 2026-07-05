package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/xray"
)

type ClientModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	InboundID uuid.UUID `gorm:"type:uuid;not null;index"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	UUID      string    `gorm:"type:varchar(36);not null;uniqueIndex"`
	Flow      string    `gorm:"type:varchar(50)"`
	Email     string    `gorm:"type:varchar(255);not null"`
	Enabled   bool      `gorm:"not null;default:true;index"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (ClientModel) TableName() string {
	return "clients"
}

func toClientModel(c *xray.Client) *ClientModel {
	return &ClientModel{
		ID:        c.ID,
		InboundID: c.InboundID,
		UserID:    c.UserID,
		UUID:      c.UUID,
		Flow:      c.Flow,
		Email:     c.Email,
		Enabled:   c.Enabled,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func toDomainClient(m *ClientModel) *xray.Client {
	return &xray.Client{
		ID:        m.ID,
		InboundID: m.InboundID,
		UserID:    m.UserID,
		UUID:      m.UUID,
		Flow:      m.Flow,
		Email:     m.Email,
		Enabled:   m.Enabled,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
