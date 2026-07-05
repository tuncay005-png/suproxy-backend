package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/usecase/auth"
	"github.com/suproxy/backend/internal/domain/user"
	"github.com/suproxy/backend/internal/infrastructure/jwt"
	"github.com/suproxy/backend/internal/interfaces/http/response"
)

type AuthHandler struct {
	registerCmd         *auth.RegisterCommand
	loginCmd            *auth.LoginCommand
	refreshTokenCmd     *auth.RefreshTokenCommand
	logoutCmd           *auth.LogoutCommand
	getCurrentUserQuery *auth.GetCurrentUserQuery
	getSessionsQuery    *auth.GetSessionsQuery
}

func NewAuthHandler(
	registerCmd *auth.RegisterCommand,
	loginCmd *auth.LoginCommand,
	refreshTokenCmd *auth.RefreshTokenCommand,
	logoutCmd *auth.LogoutCommand,
	getCurrentUserQuery *auth.GetCurrentUserQuery,
	getSessionsQuery *auth.GetSessionsQuery,
) *AuthHandler {
	return &AuthHandler{
		registerCmd:         registerCmd,
		loginCmd:            loginCmd,
		refreshTokenCmd:     refreshTokenCmd,
		logoutCmd:           logoutCmd,
		getCurrentUserQuery: getCurrentUserQuery,
		getSessionsQuery:    getSessionsQuery,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.registerCmd.Execute(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.SuccessCreated(c, result)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// Extract IP and User-Agent
	req.IPAddress = c.ClientIP()
	req.UserAgent = c.GetHeader("User-Agent")

	result, err := h.loginCmd.Execute(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.SuccessOK(c, result)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.refreshTokenCmd.Execute(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.SuccessOK(c, result)
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
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

	result, err := h.getCurrentUserQuery.Execute(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.SuccessOK(c, result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	response.SuccessNoContent(c)
}

func (h *AuthHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, user.ErrUserAlreadyExists):
		response.Conflict(c, "user already exists")
	case errors.Is(err, user.ErrUserNotFound):
		response.NotFound(c, "user not found")
	case errors.Is(err, user.ErrInvalidCredentials):
		response.Unauthorized(c, "invalid credentials")
	case errors.Is(err, user.ErrUserNotActive):
		response.Forbidden(c, "user account is not active")
	case errors.Is(err, user.ErrUserLocked):
		response.Forbidden(c, "account is temporarily locked due to multiple failed login attempts")
	case errors.Is(err, user.ErrInvalidEmail):
		response.BadRequest(c, "invalid email format")
	case errors.Is(err, jwt.ErrInvalidToken), errors.Is(err, jwt.ErrExpiredToken):
		response.Unauthorized(c, "invalid or expired token")
	default:
		response.InternalError(c, "an error occurred")
	}
}

func (h *AuthHandler) GetSessions(c *gin.Context) {
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

	sessions, err := h.getSessionsQuery.Execute(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, "failed to fetch sessions")
		return
	}

	response.SuccessOK(c, dto.ActiveSessionsResponse{
		Sessions: sessions,
		Total:    len(sessions),
	})
}

func (h *AuthHandler) LogoutSingle(c *gin.Context) {
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

	tokenID := c.Param("id")
	tokenUUID, err := uuid.Parse(tokenID)
	if err != nil {
		response.BadRequest(c, "invalid token id")
		return
	}

	err = h.logoutCmd.ExecuteSingle(
		c.Request.Context(),
		userID,
		tokenUUID,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)
	if err != nil {
		response.InternalError(c, "failed to logout")
		return
	}

	response.SuccessNoContent(c)
}

func (h *AuthHandler) LogoutAll(c *gin.Context) {
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

	err = h.logoutCmd.ExecuteAll(
		c.Request.Context(),
		userID,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)
	if err != nil {
		response.InternalError(c, "failed to logout")
		return
	}

	response.SuccessNoContent(c)
}
