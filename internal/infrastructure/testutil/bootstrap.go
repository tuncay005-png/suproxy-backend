package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suproxy/backend/internal/infrastructure/bootstrap"
	"github.com/suproxy/backend/internal/infrastructure/config"
	"github.com/suproxy/backend/internal/infrastructure/database"
	"github.com/suproxy/backend/internal/infrastructure/jwt"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

// TestApp wraps application components for integration tests
type TestApp struct {
	Config    *config.Config
	Logger    *logger.Logger
	Database  *database.Database
	JWT       *jwt.Manager
	TxManager *database.TransactionManager
	Container *bootstrap.Container
	t         *testing.T
}

// NewTestApp initializes a test application with all dependencies
func NewTestApp(t *testing.T) *TestApp {
	t.Helper()

	cfg := TestConfig()
	log := logger.New(cfg.Log.Level, cfg.Log.Format)

	// Initialize database
	db, err := database.New(cfg, log)
	require.NoError(t, err, "Failed to initialize test database")

	// Run migrations
	migrator := database.NewMigrator(cfg, log)
	require.NoError(t, migrator.Up(), "Failed to run migrations")

	// Initialize JWT manager
	jwtManager := jwt.NewManager(&cfg.JWT)

	// Initialize transaction manager
	txManager := database.NewTransactionManager(db.DB)

	// Build container
	container, err := bootstrap.BuildContainer(cfg, db.DB, log)
	require.NoError(t, err, "Failed to build dependency container")

	return &TestApp{
		Config:    cfg,
		Logger:    log,
		Database:  db,
		JWT:       jwtManager,
		TxManager: txManager,
		Container: container,
		t:         t,
	}
}

// Cleanup performs cleanup after tests
func (ta *TestApp) Cleanup() {
	if ta.Database != nil {
		_ = ta.Database.Close()
	}
}

// CleanupTables truncates all tables for test isolation
func (ta *TestApp) CleanupTables() {
	ta.t.Helper()

	testDB := &TestDatabase{
		DB:     ta.Database,
		Config: ta.Config,
		Logger: ta.Logger,
		t:      ta.t,
	}
	testDB.Cleanup()
}

