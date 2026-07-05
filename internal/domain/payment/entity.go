package payment

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	SubscriptionID uuid.UUID
	Amount         Money
	Method         Method
	Status         Status
	TransactionID  string
	Provider       string
	Metadata       map[string]string
	PaidAt         *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewPayment(userID, subscriptionID uuid.UUID, amount Money, method Method, provider string) (*Payment, error) {
	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}
	if subscriptionID == uuid.Nil {
		return nil, ErrInvalidSubscriptionID
	}

	return &Payment{
		ID:             uuid.New(),
		UserID:         userID,
		SubscriptionID: subscriptionID,
		Amount:         amount,
		Method:         method,
		Status:         StatusPending,
		Provider:       provider,
		Metadata:       make(map[string]string),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}, nil
}

func (p *Payment) MarkAsCompleted(transactionID string) error {
	if p.Status == StatusCompleted {
		return ErrPaymentAlreadyCompleted
	}
	now := time.Now().UTC()
	p.Status = StatusCompleted
	p.TransactionID = transactionID
	p.PaidAt = &now
	p.UpdatedAt = now
	return nil
}

func (p *Payment) MarkAsFailed(reason string) error {
	if p.Status == StatusFailed {
		return ErrPaymentAlreadyFailed
	}
	p.Status = StatusFailed
	p.Metadata["failure_reason"] = reason
	p.UpdatedAt = time.Now().UTC()
	return nil
}

func (p *Payment) Cancel() error {
	if p.Status == StatusCancelled {
		return ErrPaymentAlreadyCancelled
	}
	p.Status = StatusCancelled
	p.UpdatedAt = time.Now().UTC()
	return nil
}

func (p *Payment) IsCompleted() bool {
	return p.Status == StatusCompleted
}
