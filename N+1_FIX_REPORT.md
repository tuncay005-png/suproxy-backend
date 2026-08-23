# N+1 Query Problem Fix - Compact Report

## Fix Applied

**File**: `internal/application/usecase/server/list_servers_query.go`

### Before (N+1):
```go
for _, srv := range servers {
    nodeCount, _ := q.nodeRepo.CountByServerID(ctx, srv.ID) // N queries
    responses = append(responses, mapper.ToServerResponse(srv, int(nodeCount)))
}
```
**Query Count**: 1 (list servers) + N (node counts) = **21 queries** for 20 servers

### After (Batch):
```go
// Collect all server IDs
serverIDs := make([]uuid.UUID, len(servers))
for i, srv := range servers {
    serverIDs[i] = srv.ID
}

// Single batch query with GROUP BY
nodeCountMap, err := q.nodeRepo.CountByServerIDs(ctx, serverIDs) // 1 query

// Build responses
for _, srv := range servers {
    nodeCount := nodeCountMap[srv.ID]
    responses = append(responses, mapper.ToServerResponse(srv, int(nodeCount)))
}
```
**Query Count**: 1 (list servers) + 1 (batch node counts) = **2 queries total**

### New Repository Method:
```go
// internal/infrastructure/repository/node_repository.go
func (r *nodeRepository) CountByServerIDs(ctx context.Context, serverIDs []uuid.UUID) (map[uuid.UUID]int64, error) {
    // SQL: SELECT server_id, COUNT(*) FROM nodes WHERE server_id IN (...) GROUP BY server_id
}
```

## Query Reduction

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| Queries (20 servers) | 21 | 2 | **10.5x fewer** |
| Queries (100 servers) | 101 | 2 | **50.5x fewer** |

## Expected Latency (Cold Load)

**Before**:
- List servers: 815ms
- N node counts: 20 × 500ms = 10,000ms
- **Total: ~10.8s**

**After**:
- List servers: 815ms  
- Batch node counts: 50ms (single GROUP BY)
- **Total: ~865ms**

**Improvement**: 10.8s → 0.9s = **12x faster**

## Build Status

✅ `go build ./cmd/api` - SUCCESS

## Testing Required

**Backend**: Restart and check logs for timing
**Frontend**: Dashboard cold load test
**Verify**: No socket errors, no timeouts, consistent performance

## Files Changed

1. `internal/domain/node/repository.go` - Added `CountByServerIDs` interface
2. `internal/infrastructure/repository/node_repository.go` - Implemented batch method
3. `internal/application/usecase/server/list_servers_query.go` - Fixed N+1 in both Execute() and ExecuteActive()
