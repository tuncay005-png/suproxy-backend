package integration_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suproxy/backend/internal/domain/user"
	"github.com/suproxy/backend/internal/infrastructure/config"
	"github.com/suproxy/backend/internal/infrastructure/jwt"
	"github.com/suproxy/backend/internal/infrastructure/testutil"
)

func TestJWTService_GenerateAccessToken(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	cfg := &config.JWTConfig{
		SecretKey:          "test-secret-key",
		AccessTokenExpiry:  15,
		RefreshTokenExpiry: 168,
		Issuer:             "test-issuer",
	}
	jwtManager := jwt.NewManager(cfg)

	tests := []struct {
		name    string
		userID  string
		email   string
		role    string
		wantErr bool
	}{
		{"Success_User", uuid.New().String(), "user@test.com", string(user.RoleUser), false},
		{"Success_Admin", uuid.New().String(), "admin@test.com", string(user.RoleAdmin), false},
		{"Success_EmptyEmail", uuid.New().String(), "", string(user.RoleUser), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := jwtManager.GenerateAccessToken(tt.userID, tt.email, tt.role)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, token)

				// Validate the generated token
				claims, err := jwtManager.ValidateToken(token)
				require.NoError(t, err)
				assert.Equal(t, tt.userID, claims.UserID)
				assert.Equal(t, tt.email, claims.Email)
				assert.Equal(t, tt.role, claims.Role)
				assert.Equal(t, jwt.AccessToken, claims.TokenType)
			}
		})
	}
}

func TestJWTService_GenerateRefreshToken(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	cfg := &config.JWTConfig{
		SecretKey:          "test-secret-key",
		AccessTokenExpiry:  15,
		RefreshTokenExpiry: 168,
		Issuer:             "test-issuer",
	}
	jwtManager := jwt.NewManager(cfg)

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New().String()
		email := "user@test.com"
		role := string(user.RoleUser)

		token, err := jwtManager.GenerateRefreshToken(userID, email, role)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		// Validate the generated token
		claims, err := jwtManager.ValidateToken(token)
		require.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, jwt.RefreshToken, claims.TokenType)
	})
}

func TestJWTService_GenerateTokenPair(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	cfg := &config.JWTConfig{
		SecretKey:          "test-secret-key",
		AccessTokenExpiry:  15,
		RefreshTokenExpiry: 168,
		Issuer:             "test-issuer",
	}
	jwtManager := jwt.NewManager(cfg)

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New().String()
		email := "user@test.com"
		role := string(user.RoleUser)

		accessToken, refreshToken, err := jwtManager.GenerateTokenPair(userID, email, role)
		require.NoError(t, err)
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)

		// Validate access token
		accessClaims, err := jwtManager.ValidateAccessToken(accessToken)
		require.NoError(t, err)
		assert.Equal(t, userID, accessClaims.UserID)
		assert.Equal(t, jwt.AccessToken, accessClaims.TokenType)

		// Validate refresh token
		refreshClaims, err := jwtManager.ValidateRefreshToken(refreshToken)
		require.NoError(t, err)
		assert.Equal(t, userID, refreshClaims.UserID)
		assert.Equal(t, jwt.RefreshToken, refreshClaims.TokenType)
	})
}

func TestJWTService_ValidateToken(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	cfg := &config.JWTConfig{
		SecretKey:          "test-secret-key",
		AccessTokenExpiry:  15,
		RefreshTokenExpiry: 168,
		Issuer:             "test-issuer",
	}
	jwtManager := jwt.NewManager(cfg)

	t.Run("Success_ValidToken", func(t *testing.T) {
		userID := uuid.New().String()
		email := "user@test.com"
		role := string(user.RoleUser)

		token, err := jwtManager.GenerateAccessToken(userID, email, role)
		require.NoError(t, err)

		claims, err := jwtManager.ValidateToken(token)
		require.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, email, claims.Email)
		assert.Equal(t, role, claims.Role)
	})

	t.Run("Failure_InvalidToken", func(t *testing.T) {
		_, err := jwtManager.ValidateToken("invalid.token.here")
		assert.Error(t, err)
		assert.ErrorIs(t, err, jwt.ErrInvalidToken)
	})

	t.Run("Failure_EmptyToken", func(t *testing.T) {
		_, err := jwtManager.ValidateToken("")
		assert.Error(t, err)
	})

	t.Run("Failure_MalformedToken", func(t *testing.T) {
		_, err := jwtManager.ValidateToken("not-a-jwt")
		assert.Error(t, err)
	})

	t.Run("Failure_WrongSecret", func(t *testing.T) {
		userID := uuid.New().String()
		email := "user@test.com"
		role := string(user.RoleUser)

		// Generate with one secret
		token, err := jwtManager.GenerateAccessToken(userID, email, role)
		require.NoError(t, err)

		// Try to validate with different secret
		differentCfg := &config.JWTConfig{
			SecretKey:          "different-secret",
			AccessTokenExpiry:  15,
			RefreshTokenExpiry: 168,
			Issuer:             "test-issuer",
		}
		differentManager := jwt.NewManager(differentCfg)

		_, err = differentManager.ValidateToken(token)
		assert.Error(t, err)
		// The error should indicate invalid token (signature mismatch causes parse failure)
		assert.ErrorIs(t, err, jwt.ErrInvalidToken)
	})
}

func TestJWTService_ValidateAccessToken(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	cfg := &config.JWTConfig{
		SecretKey:          "test-secret-key",
		AccessTokenExpiry:  15,
		RefreshTokenExpiry: 168,
		Issuer:             "test-issuer",
	}
	jwtManager := jwt.NewManager(cfg)

	t.Run("Success_AccessToken", func(t *testing.T) {
		userID := uuid.New().String()
		token, err := jwtManager.GenerateAccessToken(userID, "user@test.com", string(user.RoleUser))
		require.NoError(t, err)

		claims, err := jwtManager.ValidateAccessToken(token)
		require.NoError(t, err)
		assert.Equal(t, jwt.AccessToken, claims.TokenType)
	})

	t.Run("Failure_RefreshToken", func(t *testing.T) {
		userID := uuid.New().String()
		token, err := jwtManager.GenerateRefreshToken(userID, "user@test.com", string(user.RoleUser))
		require.NoError(t, err)

		_, err = jwtManager.ValidateAccessToken(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected access token")
	})
}

func TestJWTService_ValidateRefreshToken(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	cfg := &config.JWTConfig{
		SecretKey:          "test-secret-key",
		AccessTokenExpiry:  15,
		RefreshTokenExpiry: 168,
		Issuer:             "test-issuer",
	}
	jwtManager := jwt.NewManager(cfg)

	t.Run("Success_RefreshToken", func(t *testing.T) {
		userID := uuid.New().String()
		token, err := jwtManager.GenerateRefreshToken(userID, "user@test.com", string(user.RoleUser))
		require.NoError(t, err)

		claims, err := jwtManager.ValidateRefreshToken(token)
		require.NoError(t, err)
		assert.Equal(t, jwt.RefreshToken, claims.TokenType)
	})

	t.Run("Failure_AccessToken", func(t *testing.T) {
		userID := uuid.New().String()
		token, err := jwtManager.GenerateAccessToken(userID, "user@test.com", string(user.RoleUser))
		require.NoError(t, err)

		_, err = jwtManager.ValidateRefreshToken(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected refresh token")
	})
}

func TestJWTService_GetUserIDFromToken(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	cfg := &config.JWTConfig{
		SecretKey:          "test-secret-key",
		AccessTokenExpiry:  15,
		RefreshTokenExpiry: 168,
		Issuer:             "test-issuer",
	}
	jwtManager := jwt.NewManager(cfg)

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New().String()
		token, err := jwtManager.GenerateAccessToken(userID, "user@test.com", string(user.RoleUser))
		require.NoError(t, err)

		extractedUserID, err := jwtManager.GetUserIDFromToken(token)
		require.NoError(t, err)
		assert.Equal(t, userID, extractedUserID)
	})

	t.Run("Failure_InvalidToken", func(t *testing.T) {
		_, err := jwtManager.GetUserIDFromToken("invalid.token")
		assert.Error(t, err)
	})
}

func TestJWTService_TokenExpiry(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	t.Run("AccessToken_ExpiryTime", func(t *testing.T) {
		cfg := &config.JWTConfig{
			SecretKey:          "test-secret-key",
			AccessTokenExpiry:  1, // 1 minute
			RefreshTokenExpiry: 168,
			Issuer:             "test-issuer",
		}
		jwtManager := jwt.NewManager(cfg)

		userID := uuid.New().String()
		token, err := jwtManager.GenerateAccessToken(userID, "user@test.com", string(user.RoleUser))
		require.NoError(t, err)

		claims, err := jwtManager.ValidateToken(token)
		require.NoError(t, err)

		// Check expiry time is approximately 1 minute from now
		expectedExpiry := time.Now().Add(1 * time.Minute)
		actualExpiry := claims.ExpiresAt.Time

		diff := actualExpiry.Sub(expectedExpiry)
		assert.Less(t, diff.Abs(), 5*time.Second) // Allow 5 second tolerance
	})

	t.Run("RefreshToken_ExpiryTime", func(t *testing.T) {
		cfg := &config.JWTConfig{
			SecretKey:          "test-secret-key",
			AccessTokenExpiry:  15,
			RefreshTokenExpiry: 1, // 1 hour
			Issuer:             "test-issuer",
		}
		jwtManager := jwt.NewManager(cfg)

		userID := uuid.New().String()
		token, err := jwtManager.GenerateRefreshToken(userID, "user@test.com", string(user.RoleUser))
		require.NoError(t, err)

		claims, err := jwtManager.ValidateToken(token)
		require.NoError(t, err)

		// Check expiry time is approximately 1 hour from now
		expectedExpiry := time.Now().Add(1 * time.Hour)
		actualExpiry := claims.ExpiresAt.Time

		diff := actualExpiry.Sub(expectedExpiry)
		assert.Less(t, diff.Abs(), 5*time.Second)
	})
}

func TestJWTService_TokenClaims(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	cfg := &config.JWTConfig{
		SecretKey:          "test-secret-key",
		AccessTokenExpiry:  15,
		RefreshTokenExpiry: 168,
		Issuer:             "test-issuer",
	}
	jwtManager := jwt.NewManager(cfg)

	t.Run("AllClaimsPresent", func(t *testing.T) {
		userID := uuid.New().String()
		email := "user@test.com"
		role := string(user.RoleUser)

		token, err := jwtManager.GenerateAccessToken(userID, email, role)
		require.NoError(t, err)

		claims, err := jwtManager.ValidateToken(token)
		require.NoError(t, err)

		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, email, claims.Email)
		assert.Equal(t, role, claims.Role)
		assert.Equal(t, cfg.Issuer, claims.Issuer)
		assert.Equal(t, userID, claims.Subject)
		assert.NotNil(t, claims.IssuedAt)
		assert.NotNil(t, claims.ExpiresAt)
		assert.NotNil(t, claims.NotBefore)
	})
}
