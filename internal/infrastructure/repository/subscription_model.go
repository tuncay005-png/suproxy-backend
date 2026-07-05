package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/subscription"
)

type SubscriptionModel struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID            uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;index"`
	PlanID            uuid.UUID `gorm:"type:uuid;not null;index"`
	Status            string    `gorm:"type:varchar(20);not null;index"`
	StartedAt         time.Time `gorm:"not null"`
	ExpiresAt         time.Time `gorm:"not null;index"`
	TrafficUsedBytes  int64     `gorm:"not null;default:0"`
	TrafficLimitBytes int64     `gorm:"not null;default:0"`
	AutoRenew         bool      `gorm:"not null;default:false"`
	CreatedAt         time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (SubscriptionModel) TableName() string {
	return "subscriptions"
}

func toSubscriptionModel(s *subscription.Subscription) *SubscriptionModel {
	return &SubscriptionModel{
		ID:                s.ID,
		UserID:            s.UserID,
		PlanID:            s.PlanID,
		Status:            string(s.Status),
		StartedAt:         s.StartedAt,
		ExpiresAt:         s.ExpiresAt,
		TrafficUsedBytes:  s.TrafficUsedBytes,
		TrafficLimitBytes: s.TrafficLimitBytes,
		AutoRenew:         s.AutoRenew,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}
}

func toDomainSubscription(m *SubscriptionModel) *subscription.Subscription {
	return &subscription.Subscription{
		ID:                m.ID,
		UserID:            m.UserID,
		PlanID:            m.PlanID,
		Status:            subscription.Status(m.Status),
		StartedAt:         m.StartedAt,
		ExpiresAt:         m.ExpiresAt,
		TrafficUsedBytes:  m.TrafficUsedBytes,
		TrafficLimitBytes: m.TrafficLimitBytes,
		AutoRenew:         m.AutoRenew,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}
