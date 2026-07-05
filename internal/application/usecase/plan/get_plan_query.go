package plan

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/subscription"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type GetPlanQuery struct {
	planRepo subscription.PlanRepository
	logger   *logger.Logger
}

func NewGetPlanQuery(planRepo subscription.PlanRepository, logger *logger.Logger) *GetPlanQuery {
	return &GetPlanQuery{
		planRepo: planRepo,
		logger:   logger,
	}
}

func (q *GetPlanQuery) Execute(ctx context.Context, planID uuid.UUID) (*dto.PlanResponse, error) {
	plan, err := q.planRepo.FindByID(ctx, planID)
	if err != nil {
		return nil, err
	}

	return mapper.ToPlanResponse(plan), nil
}
