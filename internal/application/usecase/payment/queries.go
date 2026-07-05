package payment

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
)

type GetPaymentQuery interface {
	Execute(ctx context.Context, paymentID uuid.UUID) (*dto.PaymentResponse, error)
}

type ListUserPaymentsQuery interface {
	Execute(ctx context.Context, userID uuid.UUID, from, to time.Time) (*dto.PaymentListResponse, error)
}
