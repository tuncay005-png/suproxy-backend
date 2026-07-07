package user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email Email) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*User, error)
	Count(ctx context.Context) (int64, error)
	ExistsByEmail(ctx context.Context, email Email) (bool, error)
	ListWithFilters(ctx context.Context, filters UserFilters) ([]*User, int64, error)
}

// UserFilters defines filter options for user listing
type UserFilters struct {
	Offset    int
	Limit     int
	Role      *Role
	Status    *Status
	Email     string
	SortBy    string
	SortOrder string
}
