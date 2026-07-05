package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
)

type GetServerQuery interface {
	Execute(ctx context.Context, serverID uuid.UUID) (*dto.ServerResponse, error)
}

type ListServersQuery interface {
	Execute(ctx context.Context, offset, limit int) (*dto.ServerListResponse, error)
}

type ListAvailableServersQuery interface {
	Execute(ctx context.Context) (*dto.ServerListResponse, error)
}

type ListServersByLocationQuery interface {
	Execute(ctx context.Context, countryCode string) (*dto.ServerListResponse, error)
}
