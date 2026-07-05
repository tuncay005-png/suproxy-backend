package subscription

import (
	"time"

	"github.com/google/uuid"
)

const (
	UnlimitedTraffic int64 = 0
	GBToBytes        int64 = 1024 * 1024 * 1024
)

type Plan struct {
	ID              uuid.UUID
	Name            string
	Description     string
	DurationDays    int
	TrafficLimitGB  int64  // 0 = unlimited
	DeviceLimit     int
	MaxSessions     int
	Price           Money
	Currency        string
	IsActive        bool
	Features        []string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewPlan(name, description string, durationDays int, trafficLimitGB int64, deviceLimit, maxSessions int, price Money, currency string) (*Plan, error) {
	if name == "" {
		return nil, ErrInvalidPlanName
	}
	if durationDays <= 0 {
		return nil, ErrInvalidPlanDuration
	}
	if trafficLimitGB < 0 {
		return nil, ErrInvalidTrafficLimit
	}
	if deviceLimit <= 0 {
		return nil, ErrInvalidDeviceLimit
	}
	if maxSessions <= 0 {
		return nil, ErrInvalidMaxSessions
	}

	return &Plan{
		ID:             uuid.New(),
		Name:           name,
		Description:    description,
		DurationDays:   durationDays,
		TrafficLimitGB: trafficLimitGB,
		DeviceLimit:    deviceLimit,
		MaxSessions:    maxSessions,
		Price:          price,
		Currency:       currency,
		IsActive:       true,
		Features:       []string{},
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}, nil
}

func (p *Plan) Activate() {
	p.IsActive = true
	p.UpdatedAt = time.Now().UTC()
}

func (p *Plan) Deactivate() {
	p.IsActive = false
	p.UpdatedAt = time.Now().UTC()
}

func (p *Plan) UpdatePrice(price Money, currency string) {
	p.Price = price
	p.Currency = currency
	p.UpdatedAt = time.Now().UTC()
}

func (p *Plan) UpdateDetails(description string, trafficLimitGB int64, deviceLimit, maxSessions int) error {
	if trafficLimitGB < 0 {
		return ErrInvalidTrafficLimit
	}
	if deviceLimit <= 0 {
		return ErrInvalidDeviceLimit
	}
	if maxSessions <= 0 {
		return ErrInvalidMaxSessions
	}

	p.Description = description
	p.TrafficLimitGB = trafficLimitGB
	p.DeviceLimit = deviceLimit
	p.MaxSessions = maxSessions
	p.UpdatedAt = time.Now().UTC()
	return nil
}

func (p *Plan) AddFeature(feature string) {
	for _, f := range p.Features {
		if f == feature {
			return
		}
	}
	p.Features = append(p.Features, feature)
	p.UpdatedAt = time.Now().UTC()
}

func (p *Plan) RemoveFeature(feature string) {
	newFeatures := []string{}
	for _, f := range p.Features {
		if f != feature {
			newFeatures = append(newFeatures, f)
		}
	}
	p.Features = newFeatures
	p.UpdatedAt = time.Now().UTC()
}

func (p *Plan) HasUnlimitedTraffic() bool {
	return p.TrafficLimitGB == UnlimitedTraffic
}

func (p *Plan) TrafficLimitBytes() int64 {
	if p.HasUnlimitedTraffic() {
		return 0
	}
	return p.TrafficLimitGB * GBToBytes
}

func (p *Plan) IsActiveStatus() bool {
	return p.IsActive
}
