package repository

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/subscription"
)

// StringArray is a custom type for storing []string as JSON in database
type StringArray []string

// Value implements driver.Valuer
func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

// Scan implements sql.Scanner
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(b, s)
}

type PlanModel struct {
	ID             uuid.UUID   `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name           string      `gorm:"type:varchar(100);uniqueIndex;not null"`
	Description    string      `gorm:"type:text"`
	DurationDays   int         `gorm:"not null"`
	TrafficLimitGB int64       `gorm:"not null;default:0"`
	DeviceLimit    int         `gorm:"not null"`
	MaxSessions    int         `gorm:"not null"`
	Price          int64       `gorm:"not null"`
	Currency       string      `gorm:"type:varchar(10);not null"`
	IsActive       bool        `gorm:"not null;default:true;index"`
	Features       StringArray `gorm:"type:jsonb"`
	CreatedAt      time.Time   `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time   `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (PlanModel) TableName() string {
	return "plans"
}

func toPlanModel(p *subscription.Plan) *PlanModel {
	money, _ := subscription.NewMoney(p.Price.Amount, p.Price.Currency)

	return &PlanModel{
		ID:             p.ID,
		Name:           p.Name,
		Description:    p.Description,
		DurationDays:   p.DurationDays,
		TrafficLimitGB: p.TrafficLimitGB,
		DeviceLimit:    p.DeviceLimit,
		MaxSessions:    p.MaxSessions,
		Price:          money.Amount,
		Currency:       p.Currency,
		IsActive:       p.IsActive,
		Features:       StringArray(p.Features),
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}

func toDomainPlan(m *PlanModel) (*subscription.Plan, error) {
	money, err := subscription.NewMoney(m.Price, m.Currency)
	if err != nil {
		return nil, err
	}

	return &subscription.Plan{
		ID:             m.ID,
		Name:           m.Name,
		Description:    m.Description,
		DurationDays:   m.DurationDays,
		TrafficLimitGB: m.TrafficLimitGB,
		DeviceLimit:    m.DeviceLimit,
		MaxSessions:    m.MaxSessions,
		Price:          money,
		Currency:       m.Currency,
		IsActive:       m.IsActive,
		Features:       []string(m.Features),
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}, nil
}
