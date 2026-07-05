package payment

import "errors"

var (
	ErrPaymentNotFound          = errors.New("payment not found")
	ErrPaymentAlreadyCompleted  = errors.New("payment already completed")
	ErrPaymentAlreadyFailed     = errors.New("payment already failed")
	ErrPaymentAlreadyCancelled  = errors.New("payment already cancelled")
	ErrInvalidUserID            = errors.New("invalid user id")
	ErrInvalidSubscriptionID    = errors.New("invalid subscription id")
	ErrInvalidAmount            = errors.New("invalid amount")
	ErrInvalidPaymentMethod     = errors.New("invalid payment method")
)
