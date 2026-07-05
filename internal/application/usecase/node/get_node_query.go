package node

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/node"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type GetNodeQuery struct {
	nodeRepo node.Repository
	logger   *logger.Logger
}

func NewGetNodeQuery(
	nodeRepo node.Repository,
	logger *logger.Logger,
) *GetNodeQuery {
	return &GetNodeQuery{
		nodeRepo: nodeRepo,
		logger:   logger,
	}
}

func (q *GetNodeQuery) Execute(ctx context.Context, nodeID uuid.UUID) (*dto.NodeResponse, error) {
	n, err := q.nodeRepo.FindByID(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	return mapper.ToNodeResponse(n), nil
}
