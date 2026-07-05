package node

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/node"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type ListNodesQuery struct {
	nodeRepo node.Repository
	logger   *logger.Logger
}

func NewListNodesQuery(
	nodeRepo node.Repository,
	logger *logger.Logger,
) *ListNodesQuery {
	return &ListNodesQuery{
		nodeRepo: nodeRepo,
		logger:   logger,
	}
}

func (q *ListNodesQuery) Execute(ctx context.Context, offset, limit int) (*dto.NodeListResponse, error) {
	nodes, err := q.nodeRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	total, err := q.nodeRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.NodeResponse, 0, len(nodes))
	for _, n := range nodes {
		responses = append(responses, mapper.ToNodeResponse(n))
	}

	return &dto.NodeListResponse{
		Nodes:  responses,
		Total:  total,
		Offset: offset,
		Limit:  limit,
	}, nil
}

func (q *ListNodesQuery) ExecuteByServerID(ctx context.Context, serverID uuid.UUID) (*dto.NodeListResponse, error) {
	nodes, err := q.nodeRepo.FindByServerID(ctx, serverID)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.NodeResponse, 0, len(nodes))
	for _, n := range nodes {
		responses = append(responses, mapper.ToNodeResponse(n))
	}

	return &dto.NodeListResponse{
		Nodes:  responses,
		Total:  int64(len(responses)),
		Offset: 0,
		Limit:  len(responses),
	}, nil
}
