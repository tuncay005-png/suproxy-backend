package testutil

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suproxy/backend/internal/infrastructure/config"
	"github.com/suproxy/backend/internal/infrastructure/database"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

// TestDatabase wraps a test database connection with cleanup utilities
type TestDatabase struct {
	DB     *database.Database
	Config *config.Config
	Logger *logger.Logger
	t      *testing.T
}

// NewTestDatabase creates a new test database connection
func NewTestDatabase(t *testing.T) *TestDatabase {
	t.Helper()

	cfg := TestConfig()
	log := logger.New("error", "json")

	db, err := database.New(cfg, log)
	require.NoError(t, err, "Failed to connect to test database")

	// Ping to verify connection
	require.NoError(t, db.Ping(), "Failed to ping test database")

	return &TestDatabase{
		DB:     db,
		Config: cfg,
		Logger: log,
		t:      t,
	}
}

// Close closes the database connection
func (td *TestDatabase) Close() {
	if td.DB != nil {
		_ = td.DB.Close()
	}
}

// Cleanup truncates all tables for test isolation
func (td *TestDatabase) Cleanup() {
	td.t.Helper()

	tables := []string{
		"refresh_tokens",
		"audit_logs",
		"clients",
		"reality_configs",
		"inbounds",
		"xray_instances",
		"subscriptions",
		"users",
		"nodes",
	}

	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)
		err := td.DB.DB.Exec(query).Error
		require.NoError(td.t, err, "Failed to truncate table: %s", table)
	}
}

// TruncateTable truncates a specific table
func (td *TestDatabase) TruncateTable(table string) {
	td.t.Helper()

	query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)
	err := td.DB.DB.Exec(query).Error
	require.NoError(td.t, err, "Failed to truncate table: %s", table)
}

// BeginTx starts a new transaction for testing
func (td *TestDatabase) BeginTx(ctx context.Context) error {
	return nil
}

// RunMigrations runs all database migrations
func (td *TestDatabase) RunMigrations() {
	td.t.Helper()

	migrator := database.NewMigrator(td.Config, td.Logger)
	require.NoError(td.t, migrator.Up(), "Failed to run migrations")
}

// RollbackMigrations rolls back all migrations
func (td *TestDatabase) RollbackMigrations() {
	td.t.Helper()

	migrator := database.NewMigrator(td.Config, td.Logger)
	require.NoError(td.t, migrator.Down(), "Failed to rollback migrations")
}

// ExecSQL executes raw SQL for test setup
func (td *TestDatabase) ExecSQL(query string, args ...interface{}) {
	td.t.Helper()

	err := td.DB.DB.Exec(query, args...).Error
	require.NoError(td.t, err, "Failed to execute SQL: %s", query)
}

// QueryRow executes a query that returns a single row
func (td *TestDatabase) QueryRow(query string, args ...interface{}) error {
	td.t.Helper()

	return td.DB.DB.Raw(query, args...).Row().Err()
}

// CountRows counts rows in a table
func (td *TestDatabase) CountRows(table string) int {
	td.t.Helper()

	var count int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	err := td.DB.DB.Raw(query).Count(&count).Error
	require.NoError(td.t, err, "Failed to count rows in table: %s", table)
	return int(count)
}

// TableExists checks if a table exists
func (td *TestDatabase) TableExists(table string) bool {
	td.t.Helper()

	var exists bool
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = $1
		)
	`
	err := td.DB.DB.Raw(query, table).Scan(&exists).Error
	require.NoError(td.t, err, "Failed to check table existence: %s", table)
	return exists
}
