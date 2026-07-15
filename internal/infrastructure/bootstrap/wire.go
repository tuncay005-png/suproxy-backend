package bootstrap

import (
	"github.com/gin-gonic/gin"
	adminAuditUseCasePackage "github.com/suproxy/backend/internal/application/usecase/admin/audit"
	adminClientUseCase "github.com/suproxy/backend/internal/application/usecase/admin/client"
	adminInboundUseCase "github.com/suproxy/backend/internal/application/usecase/admin/inbound"
	adminSystemUseCase "github.com/suproxy/backend/internal/application/usecase/admin/system"
	adminUserUseCase "github.com/suproxy/backend/internal/application/usecase/admin/user"
	adminXrayUseCase "github.com/suproxy/backend/internal/application/usecase/admin/xray_instance"
	authUseCase "github.com/suproxy/backend/internal/application/usecase/auth"
	nodeUseCase "github.com/suproxy/backend/internal/application/usecase/node"
	planUseCase "github.com/suproxy/backend/internal/application/usecase/plan"
	serverUseCase "github.com/suproxy/backend/internal/application/usecase/server"
	subscriptionsUseCase "github.com/suproxy/backend/internal/application/usecase/subscriptions"
	"github.com/suproxy/backend/internal/interfaces/http/handler"
	"github.com/suproxy/backend/internal/interfaces/http/router"
)

func InitializeAuthSystem(app *Application, engine *gin.Engine) {
	// Use container for dependencies
	container := app.Container

	// Initialize auth use cases
	registerCmd := authUseCase.NewRegisterCommand(
		container.UserRepository,
		container.XrayProvisioningService,
		app.Logger,
	)
	loginCmd := authUseCase.NewLoginCommand(
		container.UserRepository,
		container.RefreshTokenRepository,
		container.AuditLogRepository,
		app.JWTManager,
		app.Logger,
	)
	refreshTokenCmd := authUseCase.NewRefreshTokenCommand(
		container.UserRepository,
		container.RefreshTokenRepository,
		container.AuditLogRepository,
		app.JWTManager,
		app.Logger,
	)
	logoutCmd := authUseCase.NewLogoutCommand(
		container.RefreshTokenRepository,
		container.AuditLogRepository,
		app.Logger,
	)
	getCurrentUserQuery := authUseCase.NewGetCurrentUserQuery(container.UserRepository, app.Logger)
	getSessionsQuery := authUseCase.NewGetSessionsQuery(container.RefreshTokenRepository, app.Logger)

	// Initialize plan use cases
	listPlansQuery := planUseCase.NewListPlansQuery(container.PlanRepository, app.Logger)
	getPlanQuery := planUseCase.NewGetPlanQuery(container.PlanRepository, app.Logger)

	// Initialize subscription use cases
	getSubscriptionQuery := subscriptionsUseCase.NewGetSubscriptionQuery(
		container.SubscriptionRepository,
		container.PlanRepository,
		app.Logger,
	)

	// Initialize server use cases
	listServersQuery := serverUseCase.NewListServersQuery(
		container.ServerRepository,
		container.NodeRepository,
		app.Logger,
	)

	// Initialize node use cases
	listNodesQuery := nodeUseCase.NewListNodesQuery(container.NodeRepository, app.Logger)

	// Initialize admin user management use cases (Phase 17.2)
	listUsersQuery := adminUserUseCase.NewListUsersQuery(container.UserRepository)
	getUserQuery := adminUserUseCase.NewGetUserQuery(container.UserRepository)
	updateUserStatusCmd := adminUserUseCase.NewUpdateUserStatusCommand(
		container.UserRepository,
		container.AuditLogRepository,
	)
	updateUserRoleCmd := adminUserUseCase.NewUpdateUserRoleCommand(
		container.UserRepository,
		container.AuditLogRepository,
	)

	// Initialize admin Xray instance management use cases (Phase 17.3)
	listInstancesQuery := adminXrayUseCase.NewListInstancesQuery(container.XrayInstanceRepository)
	getInstanceQuery := adminXrayUseCase.NewGetInstanceQuery(container.XrayInstanceRepository)
	getInstanceStatsQuery := adminXrayUseCase.NewGetInstanceStatsQuery(
		container.XrayInstanceRepository,
		container.InboundRepository,
		container.ClientRepository,
	)
	startInstanceCmd := adminXrayUseCase.NewStartInstanceCommand(
		container.XrayInstanceRepository,
		container.XrayProcessManager,
		container.AuditLogRepository,
	)
	stopInstanceCmd := adminXrayUseCase.NewStopInstanceCommand(
		container.XrayInstanceRepository,
		container.XrayProcessManager,
		container.AuditLogRepository,
	)
	restartInstanceCmd := adminXrayUseCase.NewRestartInstanceCommand(
		container.XrayInstanceRepository,
		container.XrayProcessManager,
		container.AuditLogRepository,
	)
	reloadInstanceCmd := adminXrayUseCase.NewReloadInstanceCommand(
		container.XrayProvisioningService,
		container.AuditLogRepository,
	)
	checkInstanceHealthCmd := adminXrayUseCase.NewCheckInstanceHealthCommand(container.XrayProcessManager)

	// Initialize admin inbound management use cases (Phase 17.4)
	listInboundsQuery := adminInboundUseCase.NewListInboundsQuery(container.InboundRepository)
	getInboundQuery := adminInboundUseCase.NewGetInboundQuery(container.InboundRepository)
	createInboundCmd := adminInboundUseCase.NewCreateInboundCommand(
		container.InboundRepository,
		container.XrayInstanceRepository,
		container.XrayProvisioningService,
		container.AuditLogRepository,
	)
	updateInboundCmd := adminInboundUseCase.NewUpdateInboundCommand(
		container.InboundRepository,
		container.XrayProvisioningService,
		container.AuditLogRepository,
	)
	deleteInboundCmd := adminInboundUseCase.NewDeleteInboundCommand(
		container.InboundRepository,
		container.ClientRepository,
		container.XrayProvisioningService,
		container.AuditLogRepository,
	)
	enableInboundCmd := adminInboundUseCase.NewEnableInboundCommand(
		container.InboundRepository,
		container.XrayProvisioningService,
		container.AuditLogRepository,
	)
	disableInboundCmd := adminInboundUseCase.NewDisableInboundCommand(
		container.InboundRepository,
		container.XrayProvisioningService,
		container.AuditLogRepository,
	)

	// Initialize admin client management use cases (Phase 17.5)
	listClientsQuery := adminClientUseCase.NewListClientsQuery(container.ClientRepository)
	getClientQuery := adminClientUseCase.NewGetClientQuery(container.ClientRepository)
	createClientCmd := adminClientUseCase.NewCreateClientCommand(
		container.ClientRepository,
		container.InboundRepository,
		container.XrayProvisioningService,
		container.AuditLogRepository,
	)
	deleteClientCmd := adminClientUseCase.NewDeleteClientCommand(
		container.ClientRepository,
		container.InboundRepository,
		container.XrayProvisioningService,
		container.AuditLogRepository,
	)
	enableClientCmd := adminClientUseCase.NewEnableClientCommand(
		container.ClientRepository,
		container.InboundRepository,
		container.XrayProvisioningService,
		container.AuditLogRepository,
	)
	disableClientCmd := adminClientUseCase.NewDisableClientCommand(
		container.ClientRepository,
		container.InboundRepository,
		container.XrayProvisioningService,
		container.AuditLogRepository,
	)
	regenerateClientUUIDCmd := adminClientUseCase.NewRegenerateClientUUIDCommand(
		container.ClientRepository,
		container.InboundRepository,
		container.XrayProvisioningService,
		container.AuditLogRepository,
	)
	reprovisionClientCmd := adminClientUseCase.NewReprovisionClientCommand(
		container.ClientRepository,
		container.InboundRepository,
		container.XrayProvisioningService,
		container.AuditLogRepository,
	)

	// Admin Audit Log UseCases (Phase 17.6)
	listAuditLogsQuery := adminAuditUseCasePackage.NewListAuditLogsQuery(container.AuditLogRepository)
	getAuditLogQuery := adminAuditUseCasePackage.NewGetAuditLogQuery(container.AuditLogRepository)
	getAuditStatsQuery := adminAuditUseCasePackage.NewGetAuditStatsQuery(container.AuditLogRepository)

	// Admin System UseCases (Phase 17.7)
	getSystemHealthQuery := adminSystemUseCase.NewGetSystemHealthQuery(
		app.Database,
		container.XrayInstanceRepository,
		container.XrayProcessManager,
	)
	getSystemStatsQuery := adminSystemUseCase.NewGetSystemStatsQuery(
		container.UserRepository,
		container.XrayInstanceRepository,
		container.InboundRepository,
		container.ClientRepository,
		container.AuditLogRepository,
	)
	getVersionQuery := adminSystemUseCase.NewGetVersionQuery()
	getDatabaseStatusQuery := adminSystemUseCase.NewGetDatabaseStatusQuery(app.Database)
	getXraySystemStatusQuery := adminSystemUseCase.NewGetXraySystemStatusQuery(container.XrayInstanceRepository)

	// Initialize handlers
	healthHandler := handler.NewHealthHandler(app.Database)
	metricsHandler := handler.NewMetricsHandler()
	authHandler := handler.NewAuthHandler(registerCmd, loginCmd, refreshTokenCmd, logoutCmd, getCurrentUserQuery, getSessionsQuery)
	userHandler := handler.NewUserHandler()
	planHandler := handler.NewPlanHandler(listPlansQuery, getPlanQuery)
	subscriptionHandler := handler.NewSubscriptionHandler(getSubscriptionQuery)
	serverHandler := handler.NewServerHandler(listServersQuery)
	nodeHandler := handler.NewNodeHandler(listNodesQuery)
	xrayHandler := handler.NewXrayHandler()

	// Initialize admin handler
	adminHandler := handler.NewAdminHandler(
		app.Logger,
		listUsersQuery,
		getUserQuery,
		updateUserStatusCmd,
		updateUserRoleCmd,
		listInstancesQuery,
		getInstanceQuery,
		getInstanceStatsQuery,
		startInstanceCmd,
		stopInstanceCmd,
		restartInstanceCmd,
		reloadInstanceCmd,
		checkInstanceHealthCmd,
		listInboundsQuery,
		getInboundQuery,
		createInboundCmd,
		updateInboundCmd,
		deleteInboundCmd,
		enableInboundCmd,
		disableInboundCmd,
		listClientsQuery,
		getClientQuery,
		createClientCmd,
		deleteClientCmd,
		enableClientCmd,
		disableClientCmd,
		regenerateClientUUIDCmd,
		reprovisionClientCmd,
		listAuditLogsQuery,
		getAuditLogQuery,
		getAuditStatsQuery,
		getSystemHealthQuery,
		getSystemStatsQuery,
		getVersionQuery,
		getDatabaseStatusQuery,
		getXraySystemStatusQuery,
	)

	// Initialize router
	appRouter := router.NewRouter(
		engine,
		app.Logger,
		app.JWTManager,
		healthHandler,
		metricsHandler,
		authHandler,
		userHandler,
		planHandler,
		subscriptionHandler,
		serverHandler,
		nodeHandler,
		xrayHandler,
		adminHandler,
	)
	app.Router = appRouter
}
