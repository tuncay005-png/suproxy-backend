package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/suproxy/backend/internal/infrastructure/config"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type Migrator struct {
	config *config.Config
	logger *logger.Logger
}

func NewMigrator(cfg *config.Config, log *logger.Logger) *Migrator {
	return &Migrator{
		config: cfg,
		logger: log,
	}
}

func (m *Migrator) Up() error {
	mig, err := m.getInstance()
	if err != nil {
		return err
	}
	defer mig.Close()

	m.logger.Info("Running database migrations...")

	if err := mig.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	version, dirty, err := mig.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	m.logger.Info("Database migrations completed",
		"version", version,
		"dirty", dirty,
	)

	return nil
}

func (m *Migrator) Down() error {
	mig, err := m.getInstance()
	if err != nil {
		return err
	}
	defer mig.Close()

	m.logger.Info("Rolling back database migrations...")

	if err := mig.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	m.logger.Info("Database migrations rolled back")
	return nil
}

func (m *Migrator) Version() (uint, bool, error) {
	mig, err := m.getInstance()
	if err != nil {
		return 0, false, err
	}
	defer mig.Close()

	version, dirty, err := mig.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, dirty, nil
}

func (m *Migrator) getInstance() (*migrate.Migrate, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		m.config.Database.Host,
		m.config.Database.Port,
		m.config.Database.User,
		m.config.Database.Password,
		m.config.Database.DBName,
		m.config.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Get absolute path to migrations directory
	// First try from working directory
	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for migrations: %w", err)
	}
	
	// Check if migrations directory exists
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		// If not found, try from executable location (for when running tests from different directory)
		execPath, err := os.Executable()
		if err == nil {
			execDir := filepath.Dir(execPath)
			altPath := filepath.Join(execDir, "migrations")
			if _, err := os.Stat(altPath); err == nil {
				migrationsPath = altPath
			}
		}
	}

	migrationsURL := fmt.Sprintf("file://%s", filepath.ToSlash(migrationsPath))

	mig, err := migrate.NewWithDatabaseInstance(
		migrationsURL,
		"postgres",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migration instance: %w", err)
	}

	return mig, nil
}
