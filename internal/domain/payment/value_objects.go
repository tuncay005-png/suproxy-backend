package payment

import "errors"

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

// Method enum
type Method string

const (
	MethodCreditCard Method = "credit_card"
	MethodPayPal     Method = "paypal"
	MethodCrypto     Method = "crypto"
	MethodBankTransfer Method = "bank_transfer"
)

func (m Method) IsValid() bool {
	switch m {
	case MethodCreditCard, MethodPayPal, MethodCrypto, MethodBankTransfer:
		return true
	}
	return false
}

// Status enum
type Status string

const (
	StatusPending   Status = "pending"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusPending, StatusCompleted, StatusFailed, StatusCancelled:
		return true
	}
	return false
}
