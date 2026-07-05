package subscription

import "errors"

// Status represents the subscription status
type Status string

const (
	StatusActive    Status = "active"
	StatusExpired   Status = "expired"
	StatusSuspended Status = "suspended"
	StatusCancelled Status = "cancelled"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusExpired, StatusSuspended, StatusCancelled:
		return true
	default:
		return false
	}
}

func (s Status) String() string {
	return string(s)
}

// Money represents a monetary value
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
	return Money{
		Amount:   amount,
		Currency: currency,
	}, nil
}

func (m Money) IsZero() bool {
	return m.Amount == 0
}

func (m Money) Equals(other Money) bool {
	return m.Amount == other.Amount && m.Currency == other.Currency
}
