package subscription

import (
	"time"

	"github.com/google/uuid"
)

// Subscription is an aggregate root
type Subscription struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	PlanID           uuid.UUID
	Status           Status
	StartedAt        time.Time
	ExpiresAt        time.Time
	TrafficUsedBytes int64
	TrafficLimitBytes int64
	AutoRenew        bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func NewSubscription(userID, planID uuid.UUID, plan *Plan, autoRenew bool) (*Subscription, error) {
	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}
	if planID == uuid.Nil {
		return nil, ErrInvalidPlan
	}
	if plan == nil {
		return nil, ErrInvalidPlan
	}

	now := time.Now().UTC()
	expiresAt := now.AddDate(0, 0, plan.DurationDays)

	return &Subscription{
		ID:                uuid.New(),
		UserID:            userID,
		PlanID:            planID,
		Status:            StatusActive,
		StartedAt:         now,
		ExpiresAt:         expiresAt,
		TrafficUsedBytes:  0,
		TrafficLimitBytes: plan.TrafficLimitBytes(),
		AutoRenew:         autoRenew,
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}

func (s *Subscription) Activate() error {
	if s.Status == StatusActive {
		return ErrSubscriptionAlreadyActive
	}
	s.Status = StatusActive
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *Subscription) Suspend() error {
	if s.Status == StatusSuspended {
		return ErrSubscriptionAlreadySuspended
	}
	s.Status = StatusSuspended
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *Subscription) Cancel() error {
	if s.Status == StatusCancelled {
		return ErrSubscriptionAlreadyCancelled
	}
	s.Status = StatusCancelled
	s.AutoRenew = false
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *Subscription) Extend(days int) error {
	if days <= 0 {
		return ErrInvalidPlanDuration
	}
	s.ExpiresAt = s.ExpiresAt.AddDate(0, 0, days)
	
	// If subscription was expired, activate it
	if s.Status == StatusExpired {
		s.Status = StatusActive
	}
	
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *Subscription) Renew(plan *Plan) error {
	if plan == nil {
		return ErrInvalidPlan
	}

	now := time.Now().UTC()
	s.StartedAt = now
	s.ExpiresAt = now.AddDate(0, 0, plan.DurationDays)
	s.TrafficUsedBytes = 0
	s.TrafficLimitBytes = plan.TrafficLimitBytes()
	s.Status = StatusActive
	s.UpdatedAt = now
	return nil
}

func (s *Subscription) Expire() error {
	if s.Status == StatusExpired {
		return ErrSubscriptionAlreadyExpired
	}
	s.Status = StatusExpired
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *Subscription) UpdateTrafficUsage(bytesUsed int64) error {
	if bytesUsed < 0 {
		return ErrInvalidTrafficLimit
	}
	s.TrafficUsedBytes += bytesUsed
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *Subscription) ResetTrafficUsage() {
	s.TrafficUsedBytes = 0
	s.UpdatedAt = time.Now().UTC()
}

func (s *Subscription) EnableAutoRenew() {
	s.AutoRenew = true
	s.UpdatedAt = time.Now().UTC()
}

func (s *Subscription) DisableAutoRenew() {
	s.AutoRenew = false
	s.UpdatedAt = time.Now().UTC()
}

func (s *Subscription) IsActive() bool {
	return s.Status == StatusActive && !s.IsExpired()
}

func (s *Subscription) IsExpired() bool {
	return time.Now().UTC().After(s.ExpiresAt)
}

func (s *Subscription) CanConnect() bool {
	if !s.IsActive() {
		return false
	}
	
	// Check traffic limit (0 = unlimited)
	if s.TrafficLimitBytes > 0 && s.TrafficUsedBytes >= s.TrafficLimitBytes {
		return false
	}
	
	return true
}

func (s *Subscription) RemainingTraffic() int64 {
	// 0 = unlimited
	if s.TrafficLimitBytes == 0 {
		return 0
	}
	
	remaining := s.TrafficLimitBytes - s.TrafficUsedBytes
	if remaining < 0 {
		return 0
	}
	return remaining
}

func (s *Subscription) TrafficUsagePercentage() float64 {
	if s.TrafficLimitBytes == 0 {
		return 0
	}
	return float64(s.TrafficUsedBytes) / float64(s.TrafficLimitBytes) * 100
}

func (s *Subscription) DaysRemaining() int {
	if s.IsExpired() {
		return 0
	}
	duration := time.Until(s.ExpiresAt)
	return int(duration.Hours() / 24)
}

func (s *Subscription) HasUnlimitedTraffic() bool {
	return s.TrafficLimitBytes == 0
}
