package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/xray"
)

type InboundModel struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	XrayInstanceID uuid.UUID `gorm:"type:uuid;not null;index"`
	Protocol       string    `gorm:"type:varchar(20);not null"`
	Port           int       `gorm:"not null"`
	Transport      string    `gorm:"type:varchar(20);not null"`
	Security       string    `gorm:"type:varchar(20);not null"`
	Enabled        bool      `gorm:"not null;default:true;index"`
	CreatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (InboundModel) TableName() string {
	return "inbounds"
}

func toInboundModel(i *xray.Inbound) *InboundModel {
	return &InboundModel{
		ID:             i.ID,
		XrayInstanceID: i.XrayInstanceID,
		Protocol:       string(i.Protocol),
		Port:           i.Port,
		Transport:      string(i.Transport),
		Security:       string(i.Security),
		Enabled:        i.Enabled,
		CreatedAt:      i.CreatedAt,
		UpdatedAt:      i.UpdatedAt,
	}
}

func toDomainInbound(m *InboundModel) *xray.Inbound {
	return &xray.Inbound{
		ID:             m.ID,
		XrayInstanceID: m.XrayInstanceID,
		Protocol:       xray.InboundProtocol(m.Protocol),
		Port:           m.Port,
		Transport:      xray.TransportType(m.Transport),
		Security:       xray.SecurityType(m.Security),
		Enabled:        m.Enabled,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}
