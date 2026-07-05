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
		ID:           p.ID,
		Name:         p.Name,
		Description:  p.Description,
		Price:        p.Price.Amount,
		Currency:     p.Price.Currency,
		TrafficLimit: p.TrafficLimit,
		DeviceLimit:  p.DeviceLimit,
		Features:     p.Features,
		IsActive:     p.IsActive,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

func ToSubscriptionResponse(s *subscription.Subscription) *dto.SubscriptionResponse {
	if s == nil {
		return nil
	}

	return &dto.SubscriptionResponse{
		ID:            s.ID,
		UserID:        s.UserID,
		Plan:          ToPlanResponse(s.Plan),
		Status:        string(s.Status),
		Period:        string(s.Period.Duration),
		StartDate:     s.StartDate,
		EndDate:       s.EndDate,
		AutoRenew:     s.AutoRenew,
		DaysRemaining: s.DaysRemaining(),
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}
