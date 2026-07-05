package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/xray"
)

type RealityConfigModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	InboundID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	PrivateKey  string    `gorm:"type:text;not null"`
	PublicKey   string    `gorm:"type:text;not null"`
	ShortID     string    `gorm:"type:varchar(16)"`
	ServerName  string    `gorm:"type:varchar(255);not null"`
	Fingerprint string    `gorm:"type:varchar(50);not null"`
	SpiderX     string    `gorm:"type:varchar(255)"`
	CreatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (RealityConfigModel) TableName() string {
	return "reality_configs"
}

func toRealityConfigModel(r *xray.RealityConfig) *RealityConfigModel {
	return &RealityConfigModel{
		ID:          r.ID,
		InboundID:   r.InboundID,
		PrivateKey:  r.PrivateKey,
		PublicKey:   r.PublicKey,
		ShortID:     r.ShortID,
		ServerName:  r.ServerName,
		Fingerprint: r.Fingerprint,
		SpiderX:     r.SpiderX,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func toDomainRealityConfig(m *RealityConfigModel) *xray.RealityConfig {
	return &xray.RealityConfig{
		ID:          m.ID,
		InboundID:   m.InboundID,
		PrivateKey:  m.PrivateKey,
		PublicKey:   m.PublicKey,
		ShortID:     m.ShortID,
		ServerName:  m.ServerName,
		Fingerprint: m.Fingerprint,
		SpiderX:     m.SpiderX,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
