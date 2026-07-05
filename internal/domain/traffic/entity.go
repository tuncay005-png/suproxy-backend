package traffic

import (
	"time"

	"github.com/google/uuid"
)

type Usage struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	SubscriptionID uuid.UUID
	NodeID         uuid.UUID
	BytesIn        int64
	BytesOut       int64
	SessionID      uuid.UUID
	RecordedAt     time.Time
	CreatedAt      time.Time
}

func NewUsage(userID, subscriptionID, nodeID, sessionID uuid.UUID, bytesIn, bytesOut int64) (*Usage, error) {
	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}
	if subscriptionID == uuid.Nil {
		return nil, ErrInvalidSubscriptionID
	}

	return &Usage{
		ID:             uuid.New(),
		UserID:         userID,
		SubscriptionID: subscriptionID,
		NodeID:         nodeID,
		BytesIn:        bytesIn,
		BytesOut:       bytesOut,
		SessionID:      sessionID,
		RecordedAt:     time.Now().UTC(),
		CreatedAt:      time.Now().UTC(),
	}, nil
}

func (u *Usage) TotalBytes() int64 {
	return u.BytesIn + u.BytesOut
}
