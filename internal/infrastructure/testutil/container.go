package testutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestContainer provides utilities for Docker container-based testing
type TestContainer struct {
	ContainerID string
	Host        string
	Port        int
	t           *testing.T
}

// PostgresContainer represents a test PostgreSQL container
type PostgresContainer struct {
	*TestContainer
	DatabaseName string
	Username     string
	Password     string
}

// NewPostgresContainer creates a new PostgreSQL test container
// Note: This is a placeholder for future testcontainers-go integration
func NewPostgresContainer(t *testing.T) *PostgresContainer {
	t.Helper()

	// For now, use environment-configured database
	// In future, integrate with testcontainers-go
	return &PostgresContainer{
		TestContainer: &TestContainer{
			ContainerID: "",
			Host:        GetEnv("TEST_DB_HOST", "localhost"),
			Port:        GetEnvInt("TEST_DB_PORT", 5432),
			t:           t,
		},
		DatabaseName: GetEnv("TEST_DB_NAME", "suproxy_test"),
		Username:     GetEnv("TEST_DB_USER", "suproxy_test"),
		Password:     GetEnv("TEST_DB_PASSWORD", "suproxy_test"),
	}
}

// Start starts the container
func (pc *PostgresContainer) Start(ctx context.Context) error {
	pc.t.Helper()

	// Placeholder for testcontainers-go integration
	// For now, assume database is already running
	return nil
}

// Stop stops the container
func (pc *PostgresContainer) Stop(ctx context.Context) error {
	pc.t.Helper()

	// Placeholder for testcontainers-go integration
	return nil
}

// ConnectionString returns the PostgreSQL connection string
func (pc *PostgresContainer) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		pc.Host,
		pc.Port,
		pc.Username,
		pc.Password,
		pc.DatabaseName,
	)
}

// WaitForReady waits for the container to be ready
func (pc *PostgresContainer) WaitForReady(ctx context.Context, timeout time.Duration) error {
	pc.t.Helper()

	// Placeholder for readiness check
	// In real implementation, ping the database until ready
	return nil
}

// SkipIfNoDocker skips the test if Docker is not available
func SkipIfNoDocker(t *testing.T) {
	t.Helper()

	// Check if we should run integration tests
	if !IsIntegrationTest() {
		t.Skip("Skipping integration test (set INTEGRATION_TEST=true to run)")
	}
}

// RequireDocker fails the test if Docker is not available
func RequireDocker(t *testing.T) {
	t.Helper()

	// In future, check if Docker daemon is running
	// For now, just check the environment variable
	require.True(t, IsIntegrationTest(), "Docker/integration tests not enabled")
}

// ContainerCleanup ensures container cleanup
func ContainerCleanup(t *testing.T, cleanup func()) {
	t.Helper()

	t.Cleanup(cleanup)
}
