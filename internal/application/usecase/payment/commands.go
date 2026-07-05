package payment

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
)

type CreatePaymentCommand interface {
	Execute(ctx context.Context, userID uuid.UUID, req *dto.CreatePaymentRequest) (*dto.PaymentResponse, error)
}

type CompletePaymentCommand interface {
	Execute(ctx context.Context, paymentID uuid.UUID, transactionID string) error
}

type FailPaymentCommand interface {
	Execute(ctx context.Context, paymentID uuid.UUID, reason string) error
}

type CancelPaymentCommand interface {
	Execute(ctx context.Context, paymentID uuid.UUID) error
}
