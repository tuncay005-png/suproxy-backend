package server

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, server *Server) error
	FindByID(ctx context.Context, id uuid.UUID) (*Server, error)
	FindByLocation(ctx context.Context, countryCode string) ([]*Server, error)
	FindAvailable(ctx context.Context) ([]*Server, error)
	Update(ctx context.Context, server *Server) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*Server, error)
	Count(ctx context.Context) (int64, error)
	ExistsByIPAddress(ctx context.Context, ipAddress string) (bool, error)
}
