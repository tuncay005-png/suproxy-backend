package server

import (
	"context"

	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/node"
	"github.com/suproxy/backend/internal/domain/server"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type ListServersQuery struct {
	serverRepo server.Repository
	nodeRepo   node.Repository
	logger     *logger.Logger
}

func NewListServersQuery(
	serverRepo server.Repository,
	nodeRepo node.Repository,
	logger *logger.Logger,
) *ListServersQuery {
	return &ListServersQuery{
		serverRepo: serverRepo,
		nodeRepo:   nodeRepo,
		logger:     logger,
	}
}

func (q *ListServersQuery) Execute(ctx context.Context, offset, limit int) (*dto.ServerListResponse, error) {
	servers, err := q.serverRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	total, err := q.serverRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.ServerResponse, 0, len(servers))
	for _, srv := range servers {
		// Get node count for each server
		nodeCount, _ := q.nodeRepo.CountByServerID(ctx, srv.ID)
		responses = append(responses, mapper.ToServerResponse(srv, int(nodeCount)))
	}

	return &dto.ServerListResponse{
		Servers: responses,
		Total:   total,
		Offset:  offset,
		Limit:   limit,
	}, nil
}

func (q *ListServersQuery) ExecuteActive(ctx context.Context) (*dto.ServerListResponse, error) {
	servers, err := q.serverRepo.FindActive(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.ServerResponse, 0, len(servers))
	for _, srv := range servers {
		nodeCount, _ := q.nodeRepo.CountByServerID(ctx, srv.ID)
		responses = append(responses, mapper.ToServerResponse(srv, int(nodeCount)))
	}

	return &dto.ServerListResponse{
		Servers: responses,
		Total:   int64(len(responses)),
		Offset:  0,
		Limit:   len(responses),
	}, nil
}
