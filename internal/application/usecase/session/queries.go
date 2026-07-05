package session

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
)

type GetSessionQuery interface {
	Execute(ctx context.Context, sessionID uuid.UUID) (*dto.SessionResponse, error)
}

type ListUserSessionsQuery interface {
	Execute(ctx context.Context, userID uuid.UUID, from, to time.Time) (*dto.SessionListResponse, error)
}

type ListActiveSessionsQuery interface {
	Execute(ctx context.Context, userID uuid.UUID) (*dto.SessionListResponse, error)
}
