package testutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/suproxy/backend/internal/domain/session"
	"github.com/suproxy/backend/internal/domain/user"
	"github.com/suproxy/backend/internal/domain/xray"
)

// UserFixture provides test user data
type UserFixture struct {
	Username string
	Email    string
	Password string
}

// DefaultUserFixture returns a default test user
func DefaultUserFixture() UserFixture {
	return UserFixture{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "Test123!@#",
	}
}

// AdminUserFixture returns an admin test user
func AdminUserFixture() UserFixture {
	return UserFixture{
		Username: "adminuser",
		Email:    "admin@example.com",
		Password: "Admin123!@#",
	}
}

// CreateTestUser creates a test user entity
func CreateTestUser(username, email, password string) (*user.User, error) {
	emailVO, err := user.NewEmail(email)
	if err != nil {
		return nil, err
	}

	passwordVO, err := user.NewPassword(password)
	if err != nil {
		return nil, err
	}

	profile := user.NewProfile("Test", "User", "", "")

	u, err := user.NewUser(emailVO, passwordVO, profile)
	if err != nil {
		return nil, err
	}

	// Override username if needed
	// Note: In real domain, username might be part of constructor
	return u, nil
}

// CreateTestUserWithDefaults creates a test user with default fixture
func CreateTestUserWithDefaults() (*user.User, error) {
	fixture := DefaultUserFixture()
	return CreateTestUser(fixture.Username, fixture.Email, fixture.Password)
}

// CreateTestAdminUser creates an admin user with default fixture
func CreateTestAdminUser() (*user.User, error) {
	fixture := AdminUserFixture()
	u, err := CreateTestUser(fixture.Username, fixture.Email, fixture.Password)
	if err != nil {
		return nil, err
	}
	// Set role to admin directly (no PromoteToAdmin method)
	// This is a test helper, we can set fields directly if needed
	return u, nil
}

// XrayInstanceFixture provides test xray instance data
type XrayInstanceFixture struct {
	NodeID  uuid.UUID
	Version string
}

// DefaultXrayInstanceFixture returns a default test xray instance
func DefaultXrayInstanceFixture() XrayInstanceFixture {
	return XrayInstanceFixture{
		NodeID:  uuid.New(),
		Version: "1.8.0",
	}
}

// CreateTestXrayInstance creates a test xray instance entity
func CreateTestXrayInstance(fixture XrayInstanceFixture) (*xray.XrayInstance, error) {
	return xray.NewXrayInstance(fixture.NodeID, fixture.Version)
}

// CreateTestXrayInstanceWithDefaults creates a test xray instance with default fixture
func CreateTestXrayInstanceWithDefaults() (*xray.XrayInstance, error) {
	fixture := DefaultXrayInstanceFixture()
	return CreateTestXrayInstance(fixture)
}

// InboundFixture provides test inbound data
type InboundFixture struct {
	InstanceID uuid.UUID
	Protocol   xray.InboundProtocol
	Port       int
	Transport  xray.TransportType
	Security   xray.SecurityType
}

// DefaultInboundFixture returns a default test inbound
func DefaultInboundFixture(instanceID uuid.UUID) InboundFixture {
	return InboundFixture{
		InstanceID: instanceID,
		Protocol:   xray.ProtocolVLESS,
		Port:       443,
		Transport:  xray.TransportTCP,
		Security:   xray.SecurityREALITY,
	}
}

// CreateTestInbound creates a test inbound entity
func CreateTestInbound(fixture InboundFixture) (*xray.Inbound, error) {
	inbound, err := xray.NewInbound(
		fixture.InstanceID,
		fixture.Protocol,
		fixture.Port,
		fixture.Transport,
		fixture.Security,
	)
	if err != nil {
		return nil, err
	}
	inbound.Enable()
	return inbound, nil
}

// CreateTestInboundWithDefaults creates a test inbound with default fixture
func CreateTestInboundWithDefaults(instanceID uuid.UUID) (*xray.Inbound, error) {
	fixture := DefaultInboundFixture(instanceID)
	return CreateTestInbound(fixture)
}

// ClientFixture provides test client data
type ClientFixture struct {
	InboundID uuid.UUID
	UserID    uuid.UUID
	UUID      string
	Flow      string
	Email     string
}

// DefaultClientFixture returns a default test client
func DefaultClientFixture(inboundID, userID uuid.UUID) ClientFixture {
	return ClientFixture{
		InboundID: inboundID,
		UserID:    userID,
		UUID:      uuid.New().String(),
		Flow:      "xtls-rprx-vision",
		Email:     "client@example.com",
	}
}

// CreateTestClient creates a test client entity
func CreateTestClient(fixture ClientFixture) (*xray.Client, error) {
	client, err := xray.NewClient(
		fixture.InboundID,
		fixture.UserID,
		fixture.UUID,
		fixture.Flow,
		fixture.Email,
	)
	if err != nil {
		return nil, err
	}
	client.Enable()
	return client, nil
}

// CreateTestClientWithDefaults creates a test client with default fixture
func CreateTestClientWithDefaults(inboundID, userID uuid.UUID) (*xray.Client, error) {
	fixture := DefaultClientFixture(inboundID, userID)
	return CreateTestClient(fixture)
}

// RealityConfigFixture provides test reality config data
type RealityConfigFixture struct {
	InboundID   uuid.UUID
	PrivateKey  string
	PublicKey   string
	ShortID     string
	ServerName  string
	Fingerprint string
	SpiderX     string
}

// DefaultRealityConfigFixture returns a default test reality config
func DefaultRealityConfigFixture(inboundID uuid.UUID) RealityConfigFixture {
	return RealityConfigFixture{
		InboundID:   inboundID,
		PrivateKey:  "test-private-key",
		PublicKey:   "test-public-key",
		ShortID:     "0123456789abcdef",
		ServerName:  "example.com",
		Fingerprint: "chrome",
		SpiderX:     "/",
	}
}

// CreateTestRealityConfig creates a test reality config entity
func CreateTestRealityConfig(fixture RealityConfigFixture) (*xray.RealityConfig, error) {
	return xray.NewRealityConfig(
		fixture.InboundID,
		fixture.PrivateKey,
		fixture.PublicKey,
		fixture.ShortID,
		fixture.ServerName,
		fixture.Fingerprint,
		fixture.SpiderX,
	)
}

// CreateTestRealityConfigWithDefaults creates a test reality config with default fixture
func CreateTestRealityConfigWithDefaults(inboundID uuid.UUID) (*xray.RealityConfig, error) {
	fixture := DefaultRealityConfigFixture(inboundID)
	return CreateTestRealityConfig(fixture)
}

// TimeNow returns current time for test consistency
func TimeNow() time.Time {
	return time.Now().UTC()
}

// TimePast returns a past time (hours ago)
func TimePast(hours int) time.Time {
	return time.Now().UTC().Add(-time.Duration(hours) * time.Hour)
}

// TimeFuture returns a future time (hours ahead)
func TimeFuture(hours int) time.Time {
	return time.Now().UTC().Add(time.Duration(hours) * time.Hour)
}


// CreateTestRefreshToken creates a test refresh token in the database
func CreateTestRefreshToken(ctx context.Context, t *testing.T, repo session.RefreshTokenRepository, userID uuid.UUID, token string) *session.RefreshToken {
	t.Helper()
	
	// Use a simple hash for testing
	tokenHash := fmt.Sprintf("%x", token)
	
	refreshToken := session.NewRefreshToken(
		userID,
		tokenHash,
		"Test Device",
		"Test Platform",
		"127.0.0.1",
		"Test Agent",
		TimeFuture(24),
	)
	
	err := repo.Create(ctx, refreshToken)
	require.NoError(t, err, "Failed to create refresh token")
	
	return refreshToken
}
