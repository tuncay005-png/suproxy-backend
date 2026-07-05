package session

import (
	"context"

	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *RefreshToken) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	FindActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*RefreshToken, error)
	RevokeByID(ctx context.Context, id uuid.UUID) error
	RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}
