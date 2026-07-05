package mapper

import (
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/domain/subscription"
)

func ToPlanResponse(p *subscription.Plan) *dto.PlanResponse {
	if p == nil {
		return nil
	}

	return &dto.PlanResponse{
		ID:             p.ID,
		Name:           p.Name,
		Description:    p.Description,
		DurationDays:   p.DurationDays,
		TrafficLimitGB: p.TrafficLimitGB,
		DeviceLimit:    p.DeviceLimit,
		MaxSessions:    p.MaxSessions,
		Price:          p.Price.Amount,
		Currency:       p.Currency,
		IsActive:       p.IsActive,
		Features:       p.Features,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}

func ToSubscriptionResponse(s *subscription.Subscription, plan *subscription.Plan) *dto.SubscriptionResponse {
	if s == nil {
		return nil
	}

	const GBToBytes = 1024 * 1024 * 1024
	trafficUsedGB := float64(s.TrafficUsedBytes) / float64(GBToBytes)
	trafficLimitGB := float64(s.TrafficLimitBytes) / float64(GBToBytes)

	var remainingTrafficGB float64
	if s.HasUnlimitedTraffic() {
		remainingTrafficGB = -1 // Indicates unlimited
	} else {
		remainingBytes := s.RemainingTraffic()
		remainingTrafficGB = float64(remainingBytes) / float64(GBToBytes)
	}

	return &dto.SubscriptionResponse{
		ID:                  s.ID,
		UserID:              s.UserID,
		Plan:                ToPlanResponse(plan),
		Status:              string(s.Status),
		StartedAt:           s.StartedAt,
		ExpiresAt:           s.ExpiresAt,
		TrafficUsedBytes:    s.TrafficUsedBytes,
		TrafficLimitBytes:   s.TrafficLimitBytes,
		TrafficUsedGB:       trafficUsedGB,
		TrafficLimitGB:      trafficLimitGB,
		RemainingTrafficGB:  remainingTrafficGB,
		TrafficUsagePercent: s.TrafficUsagePercentage(),
		DaysRemaining:       s.DaysRemaining(),
		AutoRenew:           s.AutoRenew,
		CanConnect:          s.CanConnect(),
		HasUnlimitedTraffic: s.HasUnlimitedTraffic(),
		CreatedAt:           s.CreatedAt,
		UpdatedAt:           s.UpdatedAt,
	}
}
