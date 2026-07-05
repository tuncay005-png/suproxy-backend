package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
)

type GetUserQuery interface {
	Execute(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)
}

type GetUserByEmailQuery interface {
	Execute(ctx context.Context, email string) (*dto.UserResponse, error)
}

type ListUsersQuery interface {
	Execute(ctx context.Context, offset, limit int) (*dto.UserListResponse, error)
}
