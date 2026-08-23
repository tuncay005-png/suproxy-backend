package server

import (
"context"
"time"

"github.com/google/uuid"
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
startTotal := time.Now()

// Step 1: List servers
startList := time.Now()
servers, err := q.serverRepo.List(ctx, offset, limit)
listDuration := time.Since(startList)
if err != nil {
return nil, err
}
q.logger.Info("List servers query completed",
"duration_ms", listDuration.Milliseconds(),
"count", len(servers))

// Step 2: Count total
startCount := time.Now()
total, err := q.serverRepo.Count(ctx)
countDuration := time.Since(startCount)
if err != nil {
return nil, err
}
q.logger.Info("Count servers query completed",
"duration_ms", countDuration.Milliseconds(),
"total", total)

// Step 3: Get node counts in BATCH (fix N+1)
startNodeCounts := time.Now()

// Collect all server IDs
serverIDs := make([]uuid.UUID, len(servers))
for i, srv := range servers {
serverIDs[i] = srv.ID
}

// Single batch query instead of N queries
nodeCountMap, err := q.nodeRepo.CountByServerIDs(ctx, serverIDs)
if err != nil {
q.logger.Warn("Failed to get node counts", "error", err)
nodeCountMap = make(map[uuid.UUID]int64)
}

nodeCountsTotalDuration := time.Since(startNodeCounts)
q.logger.Info("Batch node count query completed",
"duration_ms", nodeCountsTotalDuration.Milliseconds(),
"server_count", len(servers),
"query_count", 1)

// Step 4: Build responses
responses := make([]*dto.ServerResponse, 0, len(servers))
for _, srv := range servers {
nodeCount := nodeCountMap[srv.ID]
responses = append(responses, mapper.ToServerResponse(srv, int(nodeCount)))
}

totalDuration := time.Since(startTotal)
q.logger.Info("ListServersQuery.Execute completed",
"total_duration_ms", totalDuration.Milliseconds(),
"list_duration_ms", listDuration.Milliseconds(),
"count_duration_ms", countDuration.Milliseconds(),
"node_counts_duration_ms", nodeCountsTotalDuration.Milliseconds())

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

// Batch query for active servers too
serverIDs := make([]uuid.UUID, len(servers))
for i, srv := range servers {
serverIDs[i] = srv.ID
}

nodeCountMap, err := q.nodeRepo.CountByServerIDs(ctx, serverIDs)
if err != nil {
nodeCountMap = make(map[uuid.UUID]int64)
}

responses := make([]*dto.ServerResponse, 0, len(servers))
for _, srv := range servers {
nodeCount := nodeCountMap[srv.ID]
responses = append(responses, mapper.ToServerResponse(srv, int(nodeCount)))
}

return &dto.ServerListResponse{
Servers: responses,
Total:   int64(len(responses)),
Offset:  0,
Limit:   len(responses),
}, nil
}
