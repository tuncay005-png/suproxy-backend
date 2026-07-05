package audit

import (
	"time"

	"github.com/google/uuid"
)

type Log struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Action     Action
	EntityType string
	EntityID   uuid.UUID
	IPAddress  string
	UserAgent  string
	Metadata   map[string]interface{}
	CreatedAt  time.Time
}

func NewLog(userID uuid.UUID, action Action, entityType string, entityID uuid.UUID, ipAddress, userAgent string) *Log {
	return &Log{
		ID:         uuid.New(),
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Metadata:   make(map[string]interface{}),
		CreatedAt:  time.Now().UTC(),
	}
}

func (l *Log) AddMetadata(key string, value interface{}) {
	l.Metadata[key] = value
}
