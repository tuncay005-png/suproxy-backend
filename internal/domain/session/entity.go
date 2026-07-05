package session

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	DeviceID       uuid.UUID
	NodeID         uuid.UUID
	SubscriptionID uuid.UUID
	ConnectionInfo ConnectionInfo
	Status         Status
	ConnectedAt    time.Time
	DisconnectedAt *time.Time
	Duration       int
	BytesIn        int64
	BytesOut       int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewSession(userID, deviceID, nodeID, subscriptionID uuid.UUID, connInfo ConnectionInfo) (*Session, error) {
	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}
	if deviceID == uuid.Nil {
		return nil, ErrInvalidDeviceID
	}
	if nodeID == uuid.Nil {
		return nil, ErrInvalidNodeID
	}

	return &Session{
		ID:             uuid.New(),
		UserID:         userID,
		DeviceID:       deviceID,
		NodeID:         nodeID,
		SubscriptionID: subscriptionID,
		ConnectionInfo: connInfo,
		Status:         StatusActive,
		ConnectedAt:    time.Now().UTC(),
		BytesIn:        0,
		BytesOut:       0,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}, nil
}

func (s *Session) Disconnect() error {
	if s.Status == StatusDisconnected {
		return ErrSessionAlreadyDisconnected
	}
	now := time.Now().UTC()
	s.Status = StatusDisconnected
	s.DisconnectedAt = &now
	s.Duration = int(now.Sub(s.ConnectedAt).Seconds())
	s.UpdatedAt = now
	return nil
}

func (s *Session) UpdateTraffic(bytesIn, bytesOut int64) {
	s.BytesIn = bytesIn
	s.BytesOut = bytesOut
	s.UpdatedAt = time.Now().UTC()
}

func (s *Session) TotalBytes() int64 {
	return s.BytesIn + s.BytesOut
}

func (s *Session) IsActive() bool {
	return s.Status == StatusActive
}
