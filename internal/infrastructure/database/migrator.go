package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	defer func() { _, _ = mig.Close() }()

	m.logger.Info("Running database migrations...")

	// Get current version before Up()
	currentVersion, _, _ := mig.Version()
	m.logger.Info("Current migration version before Up()", "version", currentVersion)

	if err := mig.Up(); err != nil {
		if err == migrate.ErrNoChange {
			m.logger.Info("No new migrations to apply (already up-to-date)")
		} else {
			return fmt.Errorf("failed to run migrations: %w", err)
		}
	}

	version, dirty, err := mig.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	// Debug: Log migration result
	m.logger.Info("Database migrations completed",
		"version", version,
		"dirty", dirty,
	)

	// Debug: Show version change
	if currentVersion != version {
		m.logger.Info("Migration version updated",
			"from", currentVersion,
			"to", version,
		)
	}

	return nil
}

func (m *Migrator) Down() error {
	mig, err := m.getInstance()
	if err != nil {
		return err
	}
	defer func() { _, _ = mig.Close() }()

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
	defer func() { _, _ = mig.Close() }()

	version, dirty, err := mig.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, dirty, nil
}

// CheckMigrationState queries the database to show current migration state
func (m *Migrator) CheckMigrationState() error {
	mig, err := m.getInstance()
	if err != nil {
		return err
	}
	defer func() { _, _ = mig.Close() }()

	version, dirty, err := mig.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	m.logger.Info("Current migration state",
		"version", version,
		"dirty", dirty,
	)

	return nil
}

// Force sets the migration version without running migrations
// Used to recover from dirty migration state
func (m *Migrator) Force(version int) error {
	mig, err := m.getInstance()
	if err != nil {
		return err
	}
	defer func() { _, _ = mig.Close() }()

	m.logger.Info("Forcing migration version", "version", version)

	if err := mig.Force(version); err != nil {
		return fmt.Errorf("failed to force migration version: %w", err)
	}

	m.logger.Info("Migration version forced successfully", "version", version)
	return nil
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
	// Strategy 1: Try relative to current working directory
	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for migrations: %w", err)
	}

	wd, _ := os.Getwd()
	m.logger.Info("Migration path resolution",
		"working_directory", wd,
		"relative_path", "migrations",
		"absolute_path", migrationsPath,
	)

	// Check if migrations directory exists
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		m.logger.Warn("Migrations directory not found at working directory, searching alternatives",
			"attempted_path", migrationsPath,
		)

		// Strategy 2: Try relative to executable location
		execPath, err := os.Executable()
		if err == nil {
			execDir := filepath.Dir(execPath)
			altPath := filepath.Join(execDir, "migrations")
			if _, err := os.Stat(altPath); err == nil {
				migrationsPath = altPath
				m.logger.Info("Found migrations relative to executable", "path", migrationsPath)
			} else {
				// Strategy 3: Try going up directories from working directory
				for i := 0; i < 3; i++ {
					parentMigrations := filepath.Join(wd, strings.Repeat("../", i+1), "migrations")
					absPath, err := filepath.Abs(parentMigrations)
					if err == nil {
						if _, err := os.Stat(absPath); err == nil {
							migrationsPath = absPath
							m.logger.Info("Found migrations in parent directory", "path", migrationsPath, "levels_up", i+1)
							break
						}
					}
				}
			}
		}
	}

	// Final verification
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("migrations directory not found: tried %s and parent directories", migrationsPath)
	}

	m.logger.Info("Final migrations path resolved", "path", migrationsPath)

	// List migration files for debugging
	files, err := os.ReadDir(migrationsPath)
	if err == nil {
		var migrationFiles []string
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".up.sql") {
				migrationFiles = append(migrationFiles, file.Name())
			}
		}
		m.logger.Info("Migration files discovered",
			"count", len(migrationFiles),
			"files", migrationFiles,
		)
	} else {
		m.logger.Warn("Failed to list migration files", "error", err)
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
