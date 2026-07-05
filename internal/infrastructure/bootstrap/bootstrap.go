package bootstrap

import (
	"fmt"

	"github.com/suproxy/backend/internal/infrastructure/config"
	"github.com/suproxy/backend/internal/infrastructure/database"
	"github.com/suproxy/backend/internal/infrastructure/jwt"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type Application struct {
	Config         *config.Config
	Logger         *logger.Logger
	Database       *database.Database
	JWTManager     *jwt.Manager
	TxManager      *database.TransactionManager
}

func Initialize() (*Application, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize logger
	log := logger.New(cfg.Log.Level, cfg.Log.Format)

	log.Info("Initializing SuProxy Backend",
		"version", "1.0.0",
		"environment", cfg.Environment,
	)

	// Initialize database
	db, err := database.New(cfg, log)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Test database connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Database connection verified")

	// Run migrations
	migrator := database.NewMigrator(cfg, log)
	if err := migrator.Up(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Initialize JWT manager
	jwtManager := jwt.NewManager(&cfg.JWT)

	// Initialize transaction manager
	txManager := database.NewTransactionManager(db.DB)

	log.Info("Application bootstrap completed successfully")

	return &Application{
		Config:         cfg,
		Logger:         log,
		Database:       db,
		JWTManager:     jwtManager,
		TxManager:      txManager,
	}, nil
}

func (app *Application) Shutdown() {
	app.Logger.Info("Shutting down application...")

	if err := app.Database.Close(); err != nil {
		app.Logger.Error("Failed to close database connection", "error", err)
	}

	if err := app.Logger.Sync(); err != nil {
		// Ignore sync errors on shutdown
	}

	app.Logger.Info("Application shutdown complete")
}
