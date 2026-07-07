package bootstrap

import (
	"fmt"
	"time"

	"github.com/suproxy/backend/internal/infrastructure/config"
	"github.com/suproxy/backend/internal/infrastructure/database"
	"github.com/suproxy/backend/internal/infrastructure/jwt"
	"github.com/suproxy/backend/internal/infrastructure/logger"
	"github.com/suproxy/backend/internal/infrastructure/metrics"
)

type Application struct {
	Config           *config.Config
	Logger           *logger.Logger
	Database         *database.Database
	JWTManager       *jwt.Manager
	TxManager        *database.TransactionManager
	Container        *Container
	Router           interface{ Setup() }
	MetricsCollector *metrics.Collector
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

	// Initialize Prometheus metrics
	metrics.Initialize()
	log.Info("Prometheus metrics initialized")

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

	// Build dependency container
	container, err := BuildContainer(cfg, db.DB, log)
	if err != nil {
		return nil, fmt.Errorf("failed to build dependency container: %w", err)
	}

	// Initialize metrics collector
	collectionInterval := 30 * time.Second // Default 30 seconds
	if cfg.Metrics.CollectionInterval > 0 {
		collectionInterval = time.Duration(cfg.Metrics.CollectionInterval) * time.Second
	}
	
	metricsCollector := metrics.NewCollector(
		container.UserRepository,
		container.XrayInstanceRepository,
		container.InboundRepository,
		container.ClientRepository,
		db,
		log,
		collectionInterval,
	)

	// Start metrics collector in background
	go metricsCollector.Start()
	log.Info("Metrics collector started", "interval", collectionInterval)

	log.Info("Application bootstrap completed successfully")

	return &Application{
		Config:           cfg,
		Logger:           log,
		Database:         db,
		JWTManager:       jwtManager,
		TxManager:        txManager,
		Container:        container,
		MetricsCollector: metricsCollector,
	}, nil
}

func (app *Application) Shutdown() {
	app.Logger.Info("Shutting down application...")

	// Stop metrics collector
	if app.MetricsCollector != nil {
		app.MetricsCollector.Stop()
	}

	if err := app.Database.Close(); err != nil {
		app.Logger.Error("Failed to close database connection", "error", err)
	}

	if err := app.Logger.Sync(); err != nil {
		// Ignore sync errors on shutdown
	}

	app.Logger.Info("Application shutdown complete")
}
