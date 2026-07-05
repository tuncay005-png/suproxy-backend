package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/suproxy/backend/internal/infrastructure/config"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token has expired")
	ErrInvalidSignature = errors.New("invalid token signature")
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type Claims struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

type Manager struct {
	config *config.JWTConfig
}

func NewManager(cfg *config.JWTConfig) *Manager {
	return &Manager{config: cfg}
}

// GenerateAccessToken generates a new access token
func (m *Manager) GenerateAccessToken(userID, email, role string) (string, error) {
	return m.generateToken(userID, email, role, AccessToken, time.Duration(m.config.AccessTokenExpiry)*time.Minute)
}

// GenerateRefreshToken generates a new refresh token
func (m *Manager) GenerateRefreshToken(userID, email, role string) (string, error) {
	return m.generateToken(userID, email, role, RefreshToken, time.Duration(m.config.RefreshTokenExpiry)*time.Hour)
}

// GenerateTokenPair generates both access and refresh tokens
func (m *Manager) GenerateTokenPair(userID, email, role string) (accessToken, refreshToken string, err error) {
	accessToken, err = m.GenerateAccessToken(userID, email, role)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err = m.GenerateRefreshToken(userID, email, role)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (m *Manager) generateToken(userID, email, role string, tokenType TokenType, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.config.Issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(m.config.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// ValidateToken validates and parses a JWT token
func (m *Manager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSignature
		}
		return []byte(m.config.SecretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ValidateAccessToken validates an access token
func (m *Manager) ValidateAccessToken(tokenString string) (*Claims, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != AccessToken {
		return nil, errors.New("invalid token type: expected access token")
	}

	return claims, nil
}

// ValidateRefreshToken validates a refresh token
func (m *Manager) ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != RefreshToken {
		return nil, errors.New("invalid token type: expected refresh token")
	}

	return claims, nil
}

// GetUserIDFromToken extracts user ID from token without full validation
func (m *Manager) GetUserIDFromToken(tokenString string) (string, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}
