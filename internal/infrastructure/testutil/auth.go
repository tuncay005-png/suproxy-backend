package testutil

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/suproxy/backend/internal/domain/user"
	"github.com/suproxy/backend/internal/infrastructure/jwt"
)

// AuthHelper provides authentication utilities for testing
type AuthHelper struct {
	jwtManager *jwt.Manager
	t          *testing.T
}

// NewAuthHelper creates a new auth helper
func NewAuthHelper(jwtManager *jwt.Manager, t *testing.T) *AuthHelper {
	return &AuthHelper{
		jwtManager: jwtManager,
		t:          t,
	}
}

// GenerateAccessToken generates a test access token
func (ah *AuthHelper) GenerateAccessToken(userID uuid.UUID, role user.Role) string {
	ah.t.Helper()

	token, err := ah.jwtManager.GenerateAccessToken(userID.String(), "test@example.com", string(role))
	require.NoError(ah.t, err, "Failed to generate access token")
	return token
}

// GenerateRefreshToken generates a test refresh token
func (ah *AuthHelper) GenerateRefreshToken(userID uuid.UUID) string {
	ah.t.Helper()

	token, err := ah.jwtManager.GenerateRefreshToken(userID.String(), "test@example.com", string(user.RoleUser))
	require.NoError(ah.t, err, "Failed to generate refresh token")
	return token
}

// GenerateUserTokens generates both access and refresh tokens for a user
func (ah *AuthHelper) GenerateUserTokens(userID uuid.UUID, role user.Role) (accessToken, refreshToken string) {
	ah.t.Helper()

	accessToken = ah.GenerateAccessToken(userID, role)
	refreshToken = ah.GenerateRefreshToken(userID)
	return
}

// GenerateAdminToken generates an access token for an admin user
func (ah *AuthHelper) GenerateAdminToken(userID uuid.UUID) string {
	return ah.GenerateAccessToken(userID, user.RoleAdmin)
}

// GenerateUserToken generates an access token for a regular user
func (ah *AuthHelper) GenerateUserToken(userID uuid.UUID) string {
	return ah.GenerateAccessToken(userID, user.RoleUser)
}

// ValidateToken validates a token and returns the claims
func (ah *AuthHelper) ValidateToken(token string) (*jwt.Claims, error) {
	return ah.jwtManager.ValidateToken(token)
}

// CreateAuthenticatedUser creates a user and returns tokens
func (ah *AuthHelper) CreateAuthenticatedUser(userRepo user.Repository) (u *user.User, accessToken, refreshToken string) {
	ah.t.Helper()

	ctx := context.Background()

	// Create user
	u, err := CreateTestUserWithDefaults()
	require.NoError(ah.t, err, "Failed to create test user")

	// Save to repository
	err = userRepo.Create(ctx, u)
	require.NoError(ah.t, err, "Failed to save test user")

	// Generate tokens
	accessToken, refreshToken = ah.GenerateUserTokens(u.ID, u.Role)
	return
}

// CreateAuthenticatedAdmin creates an admin user and returns tokens
func (ah *AuthHelper) CreateAuthenticatedAdmin(userRepo user.Repository) (u *user.User, accessToken, refreshToken string) {
	ah.t.Helper()

	ctx := context.Background()

	// Create admin user
	u, err := CreateTestAdminUser()
	require.NoError(ah.t, err, "Failed to create test admin user")

	// Save to repository
	err = userRepo.Create(ctx, u)
	require.NoError(ah.t, err, "Failed to save test admin user")

	// Generate tokens
	accessToken, refreshToken = ah.GenerateUserTokens(u.ID, u.Role)
	return
}

// CreateMultipleUsers creates multiple test users
func (ah *AuthHelper) CreateMultipleUsers(userRepo user.Repository, count int) []*user.User {
	ah.t.Helper()

	ctx := context.Background()
	users := make([]*user.User, count)

	for i := 0; i < count; i++ {
		uniqueID := uuid.New().String()[:8]
		fixture := UserFixture{
			Username: fmt.Sprintf("testuser_%s", uniqueID),
			Email:    fmt.Sprintf("user_%s@example.com", uniqueID),
			Password: "Test123!@#",
		}

		u, err := CreateTestUser(fixture.Username, fixture.Email, fixture.Password)
		require.NoError(ah.t, err, "Failed to create test user %d", i)

		err = userRepo.Create(ctx, u)
		require.NoError(ah.t, err, "Failed to save test user %d", i)

		users[i] = u
	}

	return users
}

// TokenPair represents an access and refresh token pair
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// GenerateTokenPair generates a token pair for testing
func (ah *AuthHelper) GenerateTokenPair(userID uuid.UUID, role user.Role) TokenPair {
	accessToken, refreshToken := ah.GenerateUserTokens(userID, role)
	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

// ExpiredToken generates an expired token for testing
func (ah *AuthHelper) ExpiredToken() string {
	// This would require modifying the JWT manager to accept custom expiry
	// For now, return an invalid token
	return "expired.token.here"
}

// InvalidToken generates an invalid token for testing
func (ah *AuthHelper) InvalidToken() string {
	return "invalid.jwt.token"
}

// MalformedToken generates a malformed token for testing
func (ah *AuthHelper) MalformedToken() string {
	return "malformed-token"
}
