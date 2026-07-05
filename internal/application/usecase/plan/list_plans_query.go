package plan

import (
	"context"

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
	plans, err := q.planRepo.List(ctx)
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
