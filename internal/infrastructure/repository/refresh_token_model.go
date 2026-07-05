package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/session"
)

type RefreshTokenModel struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null;index"`
	TokenHash  string     `gorm:"type:varchar(255);uniqueIndex;not null"`
	DeviceName string     `gorm:"type:varchar(100)"`
	Platform   string     `gorm:"type:varchar(50)"`
	IPAddress  string     `gorm:"type:varchar(45)"`
	UserAgent  string     `gorm:"type:text"`
	ExpiresAt  time.Time  `gorm:"not null;index"`
	LastUsedAt *time.Time `gorm:"type:timestamp"`
	IsRevoked  bool       `gorm:"not null;default:false;index"`
	RevokedAt  *time.Time `gorm:"type:timestamp"`
	CreatedAt  time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (RefreshTokenModel) TableName() string {
	return "refresh_tokens"
}

func toRefreshTokenModel(rt *session.RefreshToken) *RefreshTokenModel {
	return &RefreshTokenModel{
		ID:         rt.ID,
		UserID:     rt.UserID,
		TokenHash:  rt.TokenHash,
		DeviceName: rt.DeviceName,
		Platform:   rt.Platform,
		IPAddress:  rt.IPAddress,
		UserAgent:  rt.UserAgent,
		ExpiresAt:  rt.ExpiresAt,
		LastUsedAt: rt.LastUsedAt,
		IsRevoked:  rt.IsRevoked,
		RevokedAt:  rt.RevokedAt,
		CreatedAt:  rt.CreatedAt,
	}
}

func toDomainRefreshToken(m *RefreshTokenModel) *session.RefreshToken {
	return &session.RefreshToken{
		ID:         m.ID,
		UserID:     m.UserID,
		TokenHash:  m.TokenHash,
		DeviceName: m.DeviceName,
		Platform:   m.Platform,
		IPAddress:  m.IPAddress,
		UserAgent:  m.UserAgent,
		ExpiresAt:  m.ExpiresAt,
		LastUsedAt: m.LastUsedAt,
		IsRevoked:  m.IsRevoked,
		RevokedAt:  m.RevokedAt,
		CreatedAt:  m.CreatedAt,
	}
}
