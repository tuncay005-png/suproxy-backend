package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/user"
)

// ListUsersQuery handles listing users with filters
type ListUsersQuery struct {
	userRepo user.Repository
}

func NewListUsersQuery(userRepo user.Repository) *ListUsersQuery {
	return &ListUsersQuery{
		userRepo: userRepo,
	}
}

// ListUsersParams defines parameters for listing users
type ListUsersParams struct {
	Offset    int
	Limit     int
	Role      string
	Status    string
	Email     string
	SortBy    string
	SortOrder string
}

func (q *ListUsersQuery) Execute(ctx context.Context, params ListUsersParams) ([]*user.User, int64, error) {
	filters := user.UserFilters{
		Offset:    params.Offset,
		Limit:     params.Limit,
		Email:     params.Email,
		SortBy:    params.SortBy,
		SortOrder: params.SortOrder,
	}

	// Parse role filter
	if params.Role != "" {
		role := user.Role(params.Role)
		filters.Role = &role
	}

	// Parse status filter
	if params.Status != "" {
		status := user.Status(params.Status)
		filters.Status = &status
	}

	return q.userRepo.ListWithFilters(ctx, filters)
}

// GetUserQuery handles retrieving a single user by ID
type GetUserQuery struct {
	userRepo user.Repository
}

func NewGetUserQuery(userRepo user.Repository) *GetUserQuery {
	return &GetUserQuery{
		userRepo: userRepo,
	}
}

func (q *GetUserQuery) Execute(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	return q.userRepo.FindByID(ctx, userID)
}
