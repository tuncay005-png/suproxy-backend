package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/node"
	"github.com/suproxy/backend/internal/domain/server"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type GetServerQuery struct {
	serverRepo server.Repository
	nodeRepo   node.Repository
	logger     *logger.Logger
}

func NewGetServerQuery(
	serverRepo server.Repository,
	nodeRepo node.Repository,
	logger *logger.Logger,
) *GetServerQuery {
	return &GetServerQuery{
		serverRepo: serverRepo,
		nodeRepo:   nodeRepo,
		logger:     logger,
	}
}

func (q *GetServerQuery) Execute(ctx context.Context, serverID uuid.UUID) (*dto.ServerResponse, error) {
	srv, err := q.serverRepo.FindByID(ctx, serverID)
	if err != nil {
		return nil, err
	}

	// Get node count for this server
	nodeCount, _ := q.nodeRepo.CountByServerID(ctx, serverID)

	return mapper.ToServerResponse(srv, int(nodeCount)), nil
}
