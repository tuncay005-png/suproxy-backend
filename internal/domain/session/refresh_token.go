package session

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	TokenHash  string
	DeviceName string
	Platform   string
	IPAddress  string
	UserAgent  string
	ExpiresAt  time.Time
	LastUsedAt *time.Time
	IsRevoked  bool
	RevokedAt  *time.Time
	CreatedAt  time.Time
	FamilyID   uuid.UUID // Token family for reuse detection
}

func NewRefreshToken(userID uuid.UUID, tokenHash, deviceName, platform, ipAddress, userAgent string, expiresAt time.Time) *RefreshToken {
	id := uuid.New()
	return &RefreshToken{
		ID:         id,
		UserID:     userID,
		TokenHash:  tokenHash,
		DeviceName: deviceName,
		Platform:   platform,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		ExpiresAt:  expiresAt,
		IsRevoked:  false,
		CreatedAt:  time.Now().UTC(),
		FamilyID:   id, // Initially token is its own family root
	}
}

// NewRefreshTokenInFamily creates a token in an existing family (for rotation)
func NewRefreshTokenInFamily(userID, familyID uuid.UUID, tokenHash, deviceName, platform, ipAddress, userAgent string, expiresAt time.Time) *RefreshToken {
	return &RefreshToken{
		ID:         uuid.New(),
		UserID:     userID,
		TokenHash:  tokenHash,
		DeviceName: deviceName,
		Platform:   platform,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		ExpiresAt:  expiresAt,
		IsRevoked:  false,
		CreatedAt:  time.Now().UTC(),
		FamilyID:   familyID, // Inherits family from parent
	}
}

func (rt *RefreshToken) Revoke() {
	now := time.Now().UTC()
	rt.IsRevoked = true
	rt.RevokedAt = &now
}

func (rt *RefreshToken) UpdateLastUsed() {
	now := time.Now().UTC()
	rt.LastUsedAt = &now
}

func (rt *RefreshToken) IsExpired() bool {
	return time.Now().UTC().After(rt.ExpiresAt)
}

func (rt *RefreshToken) IsValid() bool {
	return !rt.IsRevoked && !rt.IsExpired()
}
