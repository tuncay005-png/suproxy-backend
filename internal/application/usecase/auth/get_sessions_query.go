package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/domain/session"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type GetSessionsQuery struct {
	refreshTokenRepo session.RefreshTokenRepository
	logger           *logger.Logger
}

func NewGetSessionsQuery(refreshTokenRepo session.RefreshTokenRepository, logger *logger.Logger) *GetSessionsQuery {
	return &GetSessionsQuery{
		refreshTokenRepo: refreshTokenRepo,
		logger:           logger,
	}
}

func (q *GetSessionsQuery) Execute(ctx context.Context, userID uuid.UUID) ([]*dto.SessionInfo, error) {
	tokens, err := q.refreshTokenRepo.FindActiveByUserID(ctx, userID)
	if err != nil {
		q.logger.Error("Failed to fetch sessions", "error", err, "user_id", userID)
		return nil, err
	}

	sessions := make([]*dto.SessionInfo, 0, len(tokens))
	for _, token := range tokens {
		sessions = append(sessions, &dto.SessionInfo{
			ID:         token.ID,
			DeviceName: token.DeviceName,
			Platform:   token.Platform,
			IPAddress:  token.IPAddress,
			LastUsedAt: token.LastUsedAt,
			CreatedAt:  token.CreatedAt,
		})
	}

	return sessions, nil
}
