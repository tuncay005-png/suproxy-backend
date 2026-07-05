package subscription

import (
	"time"

	"github.com/google/uuid"
)

// Subscription is an aggregate root
type Subscription struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Plan      *Plan
	Status    Status
	Period    Period
	StartDate time.Time
	EndDate   time.Time
	AutoRenew bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewSubscription(userID uuid.UUID, plan *Plan, period Period, autoRenew bool) (*Subscription, error) {
	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}
	if plan == nil {
		return nil, ErrInvalidPlan
	}

	now := time.Now().UTC()
	return &Subscription{
		ID:        uuid.New(),
		UserID:    userID,
		Plan:      plan,
		Status:    StatusActive,
		Period:    period,
		StartDate: now,
		EndDate:   period.CalculateEndDate(now),
		AutoRenew: autoRenew,
		CreatedAt: now,
		UpdatedAt: now,
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

func (s *Subscription) Renew(period Period) error {
	if s.Status == StatusExpired {
		s.Status = StatusActive
	}
	s.StartDate = time.Now().UTC()
	s.EndDate = period.CalculateEndDate(s.StartDate)
	s.UpdatedAt = time.Now().UTC()
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

func (s *Subscription) UpgradePlan(newPlan *Plan) error {
	if newPlan == nil {
		return ErrInvalidPlan
	}
	s.Plan = newPlan
	s.UpdatedAt = time.Now().UTC()
	return nil
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
	return s.Status == StatusActive && time.Now().UTC().Before(s.EndDate)
}

func (s *Subscription) IsExpired() bool {
	return time.Now().UTC().After(s.EndDate)
}

func (s *Subscription) DaysRemaining() int {
	if s.IsExpired() {
		return 0
	}
	duration := time.Until(s.EndDate)
	return int(duration.Hours() / 24)
}
