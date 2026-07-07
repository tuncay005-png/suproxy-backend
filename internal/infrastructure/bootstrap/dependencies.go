package bootstrap

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/suproxy/backend/internal/application/service"
	"github.com/suproxy/backend/internal/infrastructure/config"
	"github.com/suproxy/backend/internal/infrastructure/logger"
	"github.com/suproxy/backend/internal/infrastructure/repository"
	xrayImpl "github.com/suproxy/backend/internal/infrastructure/xray"
	xrayBinary "github.com/suproxy/backend/internal/infrastructure/xray/binary"
	xrayConfig "github.com/suproxy/backend/internal/infrastructure/xray/config"
	xrayRuntime "github.com/suproxy/backend/internal/infrastructure/xray/runtime"
)

// BuildContainer initializes all dependencies and returns a container
func BuildContainer(cfg *config.Config, db *gorm.DB, log *logger.Logger) (*Container, error) {
	container := &Container{}

	// Initialize repositories
	if err := initializeRepositories(container, db); err != nil {
		return nil, fmt.Errorf("failed to initialize repositories: %w", err)
	}

	// Initialize Xray infrastructure based on configuration
	if err := initializeXrayInfrastructure(container, cfg, log); err != nil {
		return nil, fmt.Errorf("failed to initialize Xray infrastructure: %w", err)
	}

	// Initialize application services
	if err := initializeApplicationServices(container, log); err != nil {
		return nil, fmt.Errorf("failed to initialize application services: %w", err)
	}

	log.Info("Dependency container built successfully",
		"xray_mode", getXrayMode(cfg),
	)

	return container, nil
}

// initializeRepositories creates all repository instances
func initializeRepositories(container *Container, db *gorm.DB) error {
	container.UserRepository = repository.NewUserRepository(db)
	container.RefreshTokenRepository = repository.NewRefreshTokenRepository(db)
	container.AuditLogRepository = repository.NewAuditLogRepository(db)
	container.SubscriptionRepository = repository.NewSubscriptionRepository(db)
	container.PlanRepository = repository.NewPlanRepository(db)
	container.ServerRepository = repository.NewServerRepository(db)
	container.NodeRepository = repository.NewNodeRepository(db)
	container.XrayInstanceRepository = repository.NewXrayInstanceRepository(db)
	container.InboundRepository = repository.NewInboundRepository(db)
	container.ClientRepository = repository.NewClientRepository(db)
	container.RealityConfigRepository = repository.NewRealityConfigRepository(db)

	return nil
}

// initializeXrayInfrastructure creates Xray-related components based on configuration
func initializeXrayInfrastructure(container *Container, cfg *config.Config, log *logger.Logger) error {
	// Decision: Use Mock or Real implementation based on configuration
	useMock := cfg.Xray.UseMock || cfg.Environment == "development" || cfg.Environment == "test"

	if useMock {
		// Mock implementations for development/testing
		container.XrayProcessManager = xrayRuntime.NewMockManager()
		container.XrayConfigWriter = xrayConfig.NewMockWriter()
		container.XrayBinaryManager = xrayBinary.NewMockManager()
		log.Info("Using Mock Xray implementations", "reason", "use_mock=true or development environment")
	} else {
		// Real implementations for production
		container.XrayProcessManager = xrayRuntime.NewRealProcessManager(
			cfg.Xray.BinaryPath,
			cfg.Xray.ConfigDir,
			cfg.Xray.LogDir,
			log,
		)
		container.XrayConfigWriter = xrayConfig.NewRealWriter(
			cfg.Xray.ConfigDir,
			cfg.Xray.BackupDir,
		)
		container.XrayBinaryManager = xrayBinary.NewRealBinaryManager(
			cfg.Xray.BinaryPath,
			cfg.Xray.InstallDir,
			log,
		)
		log.Info("Using Real Xray implementations", "binary_path", cfg.Xray.BinaryPath)
	}

	// Shared components (used by both Mock and Real)
	// Generator needs repositories, so initialize after repositories are set
	container.XrayConfigGenerator = xrayConfig.NewGenerator(
		container.XrayInstanceRepository,
		container.InboundRepository,
		container.ClientRepository,
		container.RealityConfigRepository,
	)
	container.XrayConfigValidator = xrayConfig.NewValidator()

	// Initialize Xray Kernel
	container.XrayKernel = xrayImpl.NewKernel(
		container.XrayConfigGenerator,
		container.XrayConfigValidator,
		container.XrayConfigWriter,
		container.XrayProcessManager,
		container.XrayBinaryManager,
		log,
	)

	return nil
}

// getXrayMode returns a string describing the current Xray mode
func getXrayMode(cfg *config.Config) string {
	if cfg.Xray.UseMock || cfg.Environment == "development" || cfg.Environment == "test" {
		return "mock"
	}
	return "real"
}

// initializeApplicationServices initializes application layer services
func initializeApplicationServices(container *Container, log *logger.Logger) error {
	// Xray Provisioning Service
	container.XrayProvisioningService = service.NewXrayProvisioningService(
		container.XrayInstanceRepository,
		container.InboundRepository,
		container.ClientRepository,
		container.RealityConfigRepository,
		container.AuditLogRepository,
		container.XrayConfigGenerator,
		container.XrayConfigWriter,
		container.XrayProcessManager,
		container.XrayBinaryManager,
		log,
	)

	log.Info("Application services initialized successfully")
	return nil
}
