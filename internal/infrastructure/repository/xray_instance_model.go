package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/xray"
)

type XrayInstanceModel struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	NodeID        uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex"`
	Version       string     `gorm:"type:varchar(50);not null"`
	Status        string     `gorm:"type:varchar(20);not null;default:'stopped';index"`
	ConfigVersion int        `gorm:"not null;default:1"`
	StartedAt     *time.Time `gorm:"type:timestamp with time zone"`
	StoppedAt     *time.Time `gorm:"type:timestamp with time zone"`
	CreatedAt     time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (XrayInstanceModel) TableName() string {
	return "xray_instances"
}

func toXrayInstanceModel(x *xray.XrayInstance) *XrayInstanceModel {
	return &XrayInstanceModel{
		ID:            x.ID,
		NodeID:        x.NodeID,
		Version:       x.Version,
		Status:        string(x.Status),
		ConfigVersion: x.ConfigVersion,
		StartedAt:     x.StartedAt,
		StoppedAt:     x.StoppedAt,
		CreatedAt:     x.CreatedAt,
		UpdatedAt:     x.UpdatedAt,
	}
}

func toDomainXrayInstance(m *XrayInstanceModel) *xray.XrayInstance {
	return &xray.XrayInstance{
		ID:            m.ID,
		NodeID:        m.NodeID,
		Version:       m.Version,
		Status:        xray.InstanceStatus(m.Status),
		ConfigVersion: m.ConfigVersion,
		StartedAt:     m.StartedAt,
		StoppedAt:     m.StoppedAt,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}
