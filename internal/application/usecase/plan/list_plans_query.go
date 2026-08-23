package plan

import (
	"context"
	"time"

	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/subscription"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type ListPlansQuery struct {
	planRepo subscription.PlanRepository
	logger   *logger.Logger
}

func NewListPlansQuery(planRepo subscription.PlanRepository, logger *logger.Logger) *ListPlansQuery {
	return &ListPlansQuery{
		planRepo: planRepo,
		logger:   logger,
	}
}

func (q *ListPlansQuery) ExecuteAll(ctx context.Context) (*dto.PlanListResponse, error) {
	startTotal := time.Now()

	// Step 1: List all plans
	startList := time.Now()
	plans, err := q.planRepo.List(ctx)
	listDuration := time.Since(startList)
	if err != nil {
		return nil, err
	}
	q.logger.Info("List plans query completed",
		"duration_ms", listDuration.Milliseconds(),
		"count", len(plans))

	// Step 2: Map to responses
	startMap := time.Now()
	responses := make([]*dto.PlanResponse, 0, len(plans))
	for _, plan := range plans {
		responses = append(responses, mapper.ToPlanResponse(plan))
	}
	mapDuration := time.Since(startMap)
	q.logger.Info("Plan mapping completed",
		"duration_ms", mapDuration.Milliseconds())

	totalDuration := time.Since(startTotal)
	q.logger.Info("ListPlansQuery.ExecuteAll completed",
		"total_duration_ms", totalDuration.Milliseconds(),
		"list_duration_ms", listDuration.Milliseconds(),
		"map_duration_ms", mapDuration.Milliseconds())

	return &dto.PlanListResponse{
		Plans: responses,
		Total: int64(len(responses)),
	}, nil
}

func (q *ListPlansQuery) ExecuteActive(ctx context.Context) (*dto.PlanListResponse, error) {
	plans, err := q.planRepo.FindActive(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.PlanResponse, 0, len(plans))
	for _, plan := range plans {
		responses = append(responses, mapper.ToPlanResponse(plan))
	}

	return &dto.PlanListResponse{
		Plans: responses,
		Total: int64(len(responses)),
	}, nil
}
