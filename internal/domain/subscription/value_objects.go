package subscription

import (
	"errors"
	"time"
)

// Money value object
type Money struct {
	Amount   int64
	Currency string
}

func NewMoney(amount int64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("amount cannot be negative")
	}
	if currency == "" {
		return Money{}, errors.New("currency cannot be empty")
	}
	return Money{Amount: amount, Currency: currency}, nil
}

// Period value object
type Period struct {
	Duration PeriodType
}

type PeriodType string

const (
	PeriodMonthly  PeriodType = "monthly"
	PeriodQuarterly PeriodType = "quarterly"
	PeriodYearly   PeriodType = "yearly"
)

func NewPeriod(periodType PeriodType) Period {
	return Period{Duration: periodType}
}

func (p Period) CalculateEndDate(startDate time.Time) time.Time {
	switch p.Duration {
	case PeriodMonthly:
		return startDate.AddDate(0, 1, 0)
	case PeriodQuarterly:
		return startDate.AddDate(0, 3, 0)
	case PeriodYearly:
		return startDate.AddDate(1, 0, 0)
	default:
		return startDate.AddDate(0, 1, 0)
	}
}

// Status enum
type Status string

const (
	StatusActive    Status = "active"
	StatusSuspended Status = "suspended"
	StatusCancelled Status = "cancelled"
	StatusExpired   Status = "expired"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusSuspended, StatusCancelled, StatusExpired:
		return true
	}
	return false
}
