package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
)

type CreateUserCommand interface {
	Execute(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
}

type UpdateUserCommand interface {
	Execute(ctx context.Context, userID uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
}

type ChangePasswordCommand interface {
	Execute(ctx context.Context, userID uuid.UUID, req *dto.ChangePasswordRequest) error
}

type ActivateUserCommand interface {
	Execute(ctx context.Context, userID uuid.UUID) error
}

type DeactivateUserCommand interface {
	Execute(ctx context.Context, userID uuid.UUID) error
}

type SuspendUserCommand interface {
	Execute(ctx context.Context, userID uuid.UUID) error
}

type DeleteUserCommand interface {
	Execute(ctx context.Context, userID uuid.UUID) error
}
