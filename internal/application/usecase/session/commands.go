package session

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
)

type CreateSessionCommand interface {
	Execute(ctx context.Context, userID uuid.UUID, req *dto.CreateSessionRequest) (*dto.SessionResponse, error)
}

type DisconnectSessionCommand interface {
	Execute(ctx context.Context, sessionID uuid.UUID) error
}

type UpdateSessionTrafficCommand interface {
	Execute(ctx context.Context, sessionID uuid.UUID, req *dto.UpdateSessionTrafficRequest) error
}
