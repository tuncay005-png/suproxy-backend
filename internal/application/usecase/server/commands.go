package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
)

type CreateServerCommand interface {
	Execute(ctx context.Context, req *dto.CreateServerRequest) (*dto.ServerResponse, error)
}

type UpdateServerCommand interface {
	Execute(ctx context.Context, serverID uuid.UUID, req *dto.UpdateServerRequest) (*dto.ServerResponse, error)
}

type ActivateServerCommand interface {
	Execute(ctx context.Context, serverID uuid.UUID) error
}

type DeactivateServerCommand interface {
	Execute(ctx context.Context, serverID uuid.UUID) error
}

type SetServerMaintenanceCommand interface {
	Execute(ctx context.Context, serverID uuid.UUID) error
}

type DeleteServerCommand interface {
	Execute(ctx context.Context, serverID uuid.UUID) error
}
