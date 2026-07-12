package testutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/suproxy/backend/internal/domain/node"
	"github.com/suproxy/backend/internal/domain/server"
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

// DefaultUserFixture returns a default test user with unique email
func DefaultUserFixture() UserFixture {
	// Use UUID to make email unique for each test
	uniqueID := uuid.New().String()[:8]
	return UserFixture{
		Username: "testuser_" + uniqueID,
		Email:    "test_" + uniqueID + "@example.com",
		Password: "Test123!@#",
	}
}

// AdminUserFixture returns an admin test user with unique email
func AdminUserFixture() UserFixture {
	// Use UUID to make email unique for each test
	uniqueID := uuid.New().String()[:8]
	return UserFixture{
		Username: "adminuser_" + uniqueID,
		Email:    "admin_" + uniqueID + "@example.com",
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
	emailVO, err := user.NewEmail(fixture.Email)
	if err != nil {
		return nil, err
	}

	passwordVO, err := user.NewPassword(fixture.Password)
	if err != nil {
		return nil, err
	}

	profile := user.NewProfile("Admin", "User", "", "")

	u, err := user.NewUser(emailVO, passwordVO, profile)
	if err != nil {
		return nil, err
	}
	
	// Set role to admin - access private field for test purposes
	u.Role = user.RoleAdmin
	
	return u, nil
}

// XrayInstanceFixture provides test xray instance data
type XrayInstanceFixture struct {
	NodeID  uuid.UUID
	Version string
}

// DefaultXrayInstanceFixture returns a default test xray instance with given node ID
func DefaultXrayInstanceFixture(nodeID uuid.UUID) XrayInstanceFixture {
	return XrayInstanceFixture{
		NodeID:  nodeID,
		Version: "1.8.0",
	}
}

// CreateTestXrayInstance creates a test xray instance entity
func CreateTestXrayInstance(fixture XrayInstanceFixture) (*xray.XrayInstance, error) {
	return xray.NewXrayInstance(fixture.NodeID, fixture.Version)
}

// CreateTestXrayInstanceWithDefaults creates a test xray instance with default fixture
// Requires a valid node ID
func CreateTestXrayInstanceWithDefaults(nodeID uuid.UUID) (*xray.XrayInstance, error) {
	fixture := DefaultXrayInstanceFixture(nodeID)
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

// ServerFixture provides test server data
type ServerFixture struct {
	Name     string
	Country  string
	City     string
	Hostname string
	Provider string
	IPv4     string
}

// DefaultServerFixture returns a default test server with unique hostname
func DefaultServerFixture() ServerFixture {
	// Generate unique hostname to avoid constraint violations in parallel tests
	uniqueID := uuid.New().String()[:8]
	return ServerFixture{
		Name:     "Test Server",
		Country:  "US",
		City:     "New York",
		Hostname: fmt.Sprintf("test-%s.example.com", uniqueID),
		Provider: "TestProvider",
		IPv4:     "192.168.1.1",
	}
}

// CreateTestServer creates a test server entity
func CreateTestServer(fixture ServerFixture) (*server.Server, error) {
	return server.NewServer(
		fixture.Name,
		fixture.Country,
		fixture.City,
		fixture.Hostname,
		fixture.Provider,
		fixture.IPv4,
	)
}

// CreateTestServerWithDefaults creates a test server with default fixture
func CreateTestServerWithDefaults() (*server.Server, error) {
	fixture := DefaultServerFixture()
	return CreateTestServer(fixture)
}

// NodeFixture provides test node data
type NodeFixture struct {
	ServerID         uuid.UUID
	Protocol         node.Protocol
	Port             int
	MaxUsers         int
	BandwidthLimitGB int64
}

// DefaultNodeFixture returns a default test node
func DefaultNodeFixture(serverID uuid.UUID) NodeFixture {
	return NodeFixture{
		ServerID:         serverID,
		Protocol:         node.ProtocolVLESS,
		Port:             443,
		MaxUsers:         100,
		BandwidthLimitGB: 0, // Unlimited
	}
}

// CreateTestNode creates a test node entity
func CreateTestNode(fixture NodeFixture) (*node.Node, error) {
	return node.NewNode(
		fixture.ServerID,
		fixture.Protocol,
		fixture.Port,
		fixture.MaxUsers,
		fixture.BandwidthLimitGB,
	)
}

// CreateTestNodeWithDefaults creates a test node with default fixture
func CreateTestNodeWithDefaults(serverID uuid.UUID) (*node.Node, error) {
	fixture := DefaultNodeFixture(serverID)
	return CreateTestNode(fixture)
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
