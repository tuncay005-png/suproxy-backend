package xray_instance

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/xray"
)

// ListInstancesQuery handles listing Xray instances with filters
type ListInstancesQuery struct {
	instanceRepo xray.XrayInstanceRepository
}

func NewListInstancesQuery(instanceRepo xray.XrayInstanceRepository) *ListInstancesQuery {
	return &ListInstancesQuery{
		instanceRepo: instanceRepo,
	}
}

// ListInstancesParams defines parameters for listing instances
type ListInstancesParams struct {
	Offset    int
	Limit     int
	NodeID    string
	Status    string
	SortBy    string
	SortOrder string
}

func (q *ListInstancesQuery) Execute(ctx context.Context, params ListInstancesParams) ([]*xray.XrayInstance, int64, error) {
	filters := xray.XrayInstanceFilters{
		Offset:    params.Offset,
		Limit:     params.Limit,
		SortBy:    params.SortBy,
		SortOrder: params.SortOrder,
	}

	// Parse node ID filter
	if params.NodeID != "" {
		nodeID, err := uuid.Parse(params.NodeID)
		if err != nil {
			return nil, 0, err
		}
		filters.NodeID = &nodeID
	}

	// Parse status filter
	if params.Status != "" {
		status := xray.InstanceStatus(params.Status)
		filters.Status = &status
	}

	return q.instanceRepo.ListWithFilters(ctx, filters)
}

// GetInstanceQuery handles retrieving a single instance by ID
type GetInstanceQuery struct {
	instanceRepo xray.XrayInstanceRepository
}

func NewGetInstanceQuery(instanceRepo xray.XrayInstanceRepository) *GetInstanceQuery {
	return &GetInstanceQuery{
		instanceRepo: instanceRepo,
	}
}

func (q *GetInstanceQuery) Execute(ctx context.Context, instanceID uuid.UUID) (*xray.XrayInstance, error) {
	return q.instanceRepo.FindByID(ctx, instanceID)
}

// GetInstanceStatsQuery handles retrieving instance statistics
type GetInstanceStatsQuery struct {
	instanceRepo xray.XrayInstanceRepository
	inboundRepo  xray.InboundRepository
	clientRepo   xray.ClientRepository
}

func NewGetInstanceStatsQuery(
	instanceRepo xray.XrayInstanceRepository,
	inboundRepo xray.InboundRepository,
	clientRepo xray.ClientRepository,
) *GetInstanceStatsQuery {
	return &GetInstanceStatsQuery{
		instanceRepo: instanceRepo,
		inboundRepo:  inboundRepo,
		clientRepo:   clientRepo,
	}
}

// InstanceStats represents instance statistics
type InstanceStats struct {
	TotalInbounds   int
	EnabledInbounds int
	TotalClients    int
	EnabledClients  int
}

func (q *GetInstanceStatsQuery) Execute(ctx context.Context, instanceID uuid.UUID) (*InstanceStats, error) {
	// Get all inbounds for instance
	inbounds, err := q.inboundRepo.FindByInstanceID(ctx, instanceID)
	if err != nil {
		return nil, err
	}

	// Get enabled inbounds
	enabledInbounds, err := q.inboundRepo.FindEnabledByInstanceID(ctx, instanceID)
	if err != nil {
		return nil, err
	}

	// Count clients for each inbound
	totalClients := 0
	enabledClients := 0

	for _, inbound := range inbounds {
		clients, err := q.clientRepo.FindByInboundID(ctx, inbound.ID)
		if err != nil {
			return nil, err
		}
		totalClients += len(clients)

		if inbound.IsEnabled() {
			enabledClientsList, err := q.clientRepo.FindEnabledByInboundID(ctx, inbound.ID)
			if err != nil {
				return nil, err
			}
			enabledClients += len(enabledClientsList)
		}
	}

	return &InstanceStats{
		TotalInbounds:   len(inbounds),
		EnabledInbounds: len(enabledInbounds),
		TotalClients:    totalClients,
		EnabledClients:  enabledClients,
	}, nil
}
