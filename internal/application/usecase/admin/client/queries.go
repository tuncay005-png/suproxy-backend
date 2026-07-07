package client

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/xray"
)

// ListClientsQuery handles listing clients with filters
type ListClientsQuery struct {
	clientRepo xray.ClientRepository
}

func NewListClientsQuery(clientRepo xray.ClientRepository) *ListClientsQuery {
	return &ListClientsQuery{
		clientRepo: clientRepo,
	}
}

// ListClientsParams defines parameters for listing clients
type ListClientsParams struct {
	Offset    int
	Limit     int
	InboundID string
	UserID    string
	Enabled   string
	SortBy    string
	SortOrder string
}

func (q *ListClientsQuery) Execute(ctx context.Context, params ListClientsParams) ([]*xray.Client, int64, error) {
	filters := xray.ClientFilters{
		Offset:    params.Offset,
		Limit:     params.Limit,
		SortBy:    params.SortBy,
		SortOrder: params.SortOrder,
	}

	// Parse inbound ID filter
	if params.InboundID != "" {
		inboundID, err := uuid.Parse(params.InboundID)
		if err != nil {
			return nil, 0, err
		}
		filters.InboundID = &inboundID
	}

	// Parse user ID filter
	if params.UserID != "" {
		userID, err := uuid.Parse(params.UserID)
		if err != nil {
			return nil, 0, err
		}
		filters.UserID = &userID
	}

	// Parse enabled filter
	if params.Enabled != "" {
		enabled := params.Enabled == "true"
		filters.Enabled = &enabled
	}

	return q.clientRepo.ListWithFilters(ctx, filters)
}

// GetClientQuery handles retrieving a single client by ID
type GetClientQuery struct {
	clientRepo xray.ClientRepository
}

func NewGetClientQuery(clientRepo xray.ClientRepository) *GetClientQuery {
	return &GetClientQuery{
		clientRepo: clientRepo,
	}
}

func (q *GetClientQuery) Execute(ctx context.Context, clientID uuid.UUID) (*xray.Client, error) {
	return q.clientRepo.FindByID(ctx, clientID)
}
