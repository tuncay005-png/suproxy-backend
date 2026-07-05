package node

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, node *Node) error
	FindByID(ctx context.Context, id uuid.UUID) (*Node, error)
	FindByServerID(ctx context.Context, serverID uuid.UUID) ([]*Node, error)
	FindHealthyByServerID(ctx context.Context, serverID uuid.UUID) ([]*Node, error)
	FindAvailableNodes(ctx context.Context) ([]*Node, error)
	Update(ctx context.Context, node *Node) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*Node, error)
	Count(ctx context.Context) (int64, error)
	CountByServerID(ctx context.Context, serverID uuid.UUID) (int64, error)
}
