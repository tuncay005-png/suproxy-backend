package inbound

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/xray"
)

// ListInboundsQuery handles listing inbounds with filters
type ListInboundsQuery struct {
	inboundRepo xray.InboundRepository
}

func NewListInboundsQuery(inboundRepo xray.InboundRepository) *ListInboundsQuery {
	return &ListInboundsQuery{
		inboundRepo: inboundRepo,
	}
}

// ListInboundsParams defines parameters for listing inbounds
type ListInboundsParams struct {
	Offset     int
	Limit      int
	InstanceID string
	Protocol   string
	Enabled    string
	SortBy     string
	SortOrder  string
}

func (q *ListInboundsQuery) Execute(ctx context.Context, params ListInboundsParams) ([]*xray.Inbound, int64, error) {
	filters := xray.InboundFilters{
		Offset:    params.Offset,
		Limit:     params.Limit,
		SortBy:    params.SortBy,
		SortOrder: params.SortOrder,
	}

	// Parse instance ID filter
	if params.InstanceID != "" {
		instanceID, err := uuid.Parse(params.InstanceID)
		if err != nil {
			return nil, 0, err
		}
		filters.InstanceID = &instanceID
	}

	// Parse protocol filter
	if params.Protocol != "" {
		protocol := xray.InboundProtocol(params.Protocol)
		filters.Protocol = &protocol
	}

	// Parse enabled filter
	if params.Enabled != "" {
		enabled := params.Enabled == "true"
		filters.Enabled = &enabled
	}

	return q.inboundRepo.ListWithFilters(ctx, filters)
}

// GetInboundQuery handles retrieving a single inbound by ID
type GetInboundQuery struct {
	inboundRepo xray.InboundRepository
}

func NewGetInboundQuery(inboundRepo xray.InboundRepository) *GetInboundQuery {
	return &GetInboundQuery{
		inboundRepo: inboundRepo,
	}
}

func (q *GetInboundQuery) Execute(ctx context.Context, inboundID uuid.UUID) (*xray.Inbound, error) {
	return q.inboundRepo.FindByID(ctx, inboundID)
}
