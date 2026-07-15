package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/usecase/plan"
	"github.com/suproxy/backend/internal/interfaces/http/response"
)

// PlanHandler handles plan-related endpoints
type PlanHandler struct {
	listPlansQuery *plan.ListPlansQuery
	getPlanQuery   *plan.GetPlanQuery
}

// NewPlanHandler creates a new plan handler
func NewPlanHandler(
	listPlansQuery *plan.ListPlansQuery,
	getPlanQuery *plan.GetPlanQuery,
) *PlanHandler {
	return &PlanHandler{
		listPlansQuery: listPlansQuery,
		getPlanQuery:   getPlanQuery,
	}
}

// ListPlans godoc
// @Summary List all plans
// @Description Get a list of all available subscription plans
// @Tags plans
// @Accept json
// @Produce json
// @Success 200 {object} dto.PlanListResponse
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/plans [get]
func (h *PlanHandler) ListPlans(c *gin.Context) {
	plans, err := h.listPlansQuery.ExecuteAll(c.Request.Context())
	if err != nil {
		response.InternalError(c, "failed to fetch plans")
		return
	}

	response.SuccessOK(c, plans)
}

// GetPlan godoc
// @Summary Get plan by ID
// @Description Get details of a specific subscription plan
// @Tags plans
// @Accept json
// @Produce json
// @Param id path string true "Plan ID"
// @Success 200 {object} dto.PlanResponse
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/plans/{id} [get]
func (h *PlanHandler) GetPlan(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := uuid.Parse(planIDStr)
	if err != nil {
		response.BadRequest(c, "invalid plan id")
		return
	}

	plan, err := h.getPlanQuery.Execute(c.Request.Context(), planID)
	if err != nil {
		response.NotFound(c, "plan not found")
		return
	}

	response.SuccessOK(c, plan)
}
