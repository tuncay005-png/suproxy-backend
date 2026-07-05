package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/node"
)

type NodeModel struct {
	ID                  uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ServerID            uuid.UUID `gorm:"type:uuid;not null;index"`
	Protocol            string    `gorm:"type:varchar(20);not null"`
	Port                int       `gorm:"not null"`
	MaxUsers            int       `gorm:"not null"`
	CurrentUsers        int       `gorm:"not null;default:0"`
	BandwidthLimitBytes int64     `gorm:"not null;default:0"`
	BandwidthUsedBytes  int64     `gorm:"not null;default:0"`
	CPUUsage            float64   `gorm:"type:decimal(5,2);not null;default:0"`
	RAMUsage            float64   `gorm:"type:decimal(5,2);not null;default:0"`
	LatencyMs           int       `gorm:"not null;default:0"`
	Version             string    `gorm:"type:varchar(50)"`
	HealthStatus        string    `gorm:"type:varchar(20);not null;default:'unknown';index"`
	CreatedAt           time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt           time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (NodeModel) TableName() string {
	return "nodes"
}

func toNodeModel(n *node.Node) *NodeModel {
	return &NodeModel{
		ID:                  n.ID,
		ServerID:            n.ServerID,
		Protocol:            string(n.Protocol),
		Port:                n.Port,
		MaxUsers:            n.MaxUsers,
		CurrentUsers:        n.CurrentUsers,
		BandwidthLimitBytes: n.BandwidthLimitBytes,
		BandwidthUsedBytes:  n.BandwidthUsedBytes,
		CPUUsage:            n.CPUUsage,
		RAMUsage:            n.RAMUsage,
		LatencyMs:           n.LatencyMs,
		Version:             n.Version,
		HealthStatus:        string(n.HealthStatus),
		CreatedAt:           n.CreatedAt,
		UpdatedAt:           n.UpdatedAt,
	}
}

func toDomainNode(m *NodeModel) *node.Node {
	return &node.Node{
		ID:                  m.ID,
		ServerID:            m.ServerID,
		Protocol:            node.Protocol(m.Protocol),
		Port:                m.Port,
		MaxUsers:            m.MaxUsers,
		CurrentUsers:        m.CurrentUsers,
		BandwidthLimitBytes: m.BandwidthLimitBytes,
		BandwidthUsedBytes:  m.BandwidthUsedBytes,
		CPUUsage:            m.CPUUsage,
		RAMUsage:            m.RAMUsage,
		LatencyMs:           m.LatencyMs,
		Version:             m.Version,
		HealthStatus:        node.HealthStatus(m.HealthStatus),
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
	}
}
