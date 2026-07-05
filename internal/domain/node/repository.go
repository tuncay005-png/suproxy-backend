package node

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, node *Node) error
	FindByID(ctx context.Context, id uuid.UUID) (*Node, error)
	FindByServerID(ctx context.Context, serverID uuid.UUID) ([]*Node, error)
	FindActiveByServerID(ctx context.Context, serverID uuid.UUID) ([]*Node, error)
	Update(ctx context.Context, node *Node) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountByServerID(ctx context.Context, serverID uuid.UUID) (int64, error)
}
