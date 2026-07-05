package session

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, session *Session) error
	FindByID(ctx context.Context, id uuid.UUID) (*Session, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]*Session, error)
	FindActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*Session, error)
	FindActiveByDeviceID(ctx context.Context, deviceID uuid.UUID) (*Session, error)
	Update(ctx context.Context, session *Session) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountActiveByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
}
