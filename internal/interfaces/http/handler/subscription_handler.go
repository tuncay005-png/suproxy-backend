package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/usecase/subscriptions"
	"github.com/suproxy/backend/internal/infrastructure/jwt"
	"github.com/suproxy/backend/internal/interfaces/http/response"
)

// SubscriptionHandler handles subscription-related endpoints
type SubscriptionHandler struct {
	getSubscriptionQuery *subscriptions.GetSubscriptionQuery
}

// NewSubscriptionHandler creates a new subscription handler
func NewSubscriptionHandler(
	getSubscriptionQuery *subscriptions.GetSubscriptionQuery,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		getSubscriptionQuery: getSubscriptionQuery,
	}
}

// GetMySubscription godoc
// @Summary Get current user's subscription
// @Description Get the subscription details of the authenticated user
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.SubscriptionResponse
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/subscriptions/me [get]
func (h *SubscriptionHandler) GetMySubscription(c *gin.Context) {
	userIDStr, exists := jwt.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "user not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	subscription, err := h.getSubscriptionQuery.ExecuteByUserID(c.Request.Context(), userID)
	if err != nil {
		response.NotFound(c, "subscription not found")
		return
	}

	response.SuccessOK(c, subscription)
}
