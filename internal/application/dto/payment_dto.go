package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreatePaymentRequest struct {
	SubscriptionID uuid.UUID `json:"subscription_id" binding:"required"`
	Amount         int64     `json:"amount" binding:"required"`
	Currency       string    `json:"currency" binding:"required"`
	Method         string    `json:"method" binding:"required"`
	Provider       string    `json:"provider" binding:"required"`
}

type PaymentResponse struct {
	ID             uuid.UUID         `json:"id"`
	UserID         uuid.UUID         `json:"user_id"`
	SubscriptionID uuid.UUID         `json:"subscription_id"`
	Amount         int64             `json:"amount"`
	Currency       string            `json:"currency"`
	Method         string            `json:"method"`
	Status         string            `json:"status"`
	TransactionID  string            `json:"transaction_id,omitempty"`
	Provider       string            `json:"provider"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	PaidAt         *time.Time        `json:"paid_at"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

type PaymentListResponse struct {
	Payments []*PaymentResponse `json:"payments"`
	Total    int64              `json:"total"`
}
