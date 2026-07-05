package subscription

import (
	"time"

	"github.com/google/uuid"
)

type Plan struct {
	ID                uuid.UUID
	Name              string
	Description       string
	Price             Money
	TrafficLimit      int64
	DeviceLimit       int
	ServerAccessLevel int
	Features          []string
	IsActive          bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func NewPlan(name, description string, price Money, trafficLimit int64, deviceLimit int) (*Plan, error) {
	if name == "" {
		return nil, ErrInvalidPlanName
	}
	if trafficLimit <= 0 {
		return nil, ErrInvalidTrafficLimit
	}
	if deviceLimit <= 0 {
		return nil, ErrInvalidDeviceLimit
	}

	return &Plan{
		ID:                uuid.New(),
		Name:              name,
		Description:       description,
		Price:             price,
		TrafficLimit:      trafficLimit,
		DeviceLimit:       deviceLimit,
		ServerAccessLevel: 1,
		Features:          []string{},
		IsActive:          true,
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
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

func (p *Plan) UpdatePrice(price Money) {
	p.Price = price
	p.UpdatedAt = time.Now().UTC()
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
