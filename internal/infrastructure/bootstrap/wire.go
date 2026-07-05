package bootstrap

import (
	"github.com/gin-gonic/gin"
	authUseCase "github.com/suproxy/backend/internal/application/usecase/auth"
	nodeUseCase "github.com/suproxy/backend/internal/application/usecase/node"
	planUseCase "github.com/suproxy/backend/internal/application/usecase/plan"
	serverUseCase "github.com/suproxy/backend/internal/application/usecase/server"
	subscriptionsUseCase "github.com/suproxy/backend/internal/application/usecase/subscriptions"
	"github.com/suproxy/backend/internal/infrastructure/repository"
	"github.com/suproxy/backend/internal/interfaces/http/handler"
	"github.com/suproxy/backend/internal/interfaces/http/router"
)

func InitializeAuthSystem(app *Application, engine *gin.Engine) {
	// Initialize repositories
	userRepo := repository.NewUserRepository(app.Database.DB)
	refreshTokenRepo := repository.NewRefreshTokenRepository(app.Database.DB)
	auditRepo := repository.NewAuditLogRepository(app.Database.DB)
	planRepo := repository.NewPlanRepository(app.Database.DB)
	subscriptionRepo := repository.NewSubscriptionRepository(app.Database.DB)
	serverRepo := repository.NewServerRepository(app.Database.DB)
	nodeRepo := repository.NewNodeRepository(app.Database.DB)

	// Initialize auth use cases
	registerCmd := authUseCase.NewRegisterCommand(userRepo, app.Logger)
	loginCmd := authUseCase.NewLoginCommand(userRepo, refreshTokenRepo, auditRepo, app.JWTManager, app.Logger)
	refreshTokenCmd := authUseCase.NewRefreshTokenCommand(userRepo, refreshTokenRepo, auditRepo, app.JWTManager, app.Logger)
	logoutCmd := authUseCase.NewLogoutCommand(refreshTokenRepo, auditRepo, app.Logger)
	getCurrentUserQuery := authUseCase.NewGetCurrentUserQuery(userRepo, app.Logger)
	getSessionsQuery := authUseCase.NewGetSessionsQuery(refreshTokenRepo, app.Logger)

	// Initialize plan use cases
	listPlansQuery := planUseCase.NewListPlansQuery(planRepo, app.Logger)
	getPlanQuery := planUseCase.NewGetPlanQuery(planRepo, app.Logger)

	// Initialize subscription use cases
	getSubscriptionQuery := subscriptionsUseCase.NewGetSubscriptionQuery(subscriptionRepo, planRepo, app.Logger)

	// Initialize server use cases
	listServersQuery := serverUseCase.NewListServersQuery(serverRepo, nodeRepo, app.Logger)

	// Initialize node use cases
	listNodesQuery := nodeUseCase.NewListNodesQuery(nodeRepo, app.Logger)

	// Initialize handlers
	healthHandler := handler.NewHealthHandler(app.Database)
	authHandler := handler.NewAuthHandler(registerCmd, loginCmd, refreshTokenCmd, logoutCmd, getCurrentUserQuery, getSessionsQuery)
	userHandler := handler.NewUserHandler()
	planHandler := handler.NewPlanHandler(listPlansQuery, getPlanQuery)
	subscriptionHandler := handler.NewSubscriptionHandler(getSubscriptionQuery)
	serverHandler := handler.NewServerHandler(listServersQuery)
	nodeHandler := handler.NewNodeHandler(listNodesQuery)
	xrayHandler := handler.NewXrayHandler()

	// Initialize router
	appRouter := router.NewRouter(
		engine,
		app.Logger,
		app.JWTManager,
		healthHandler,
		authHandler,
		userHandler,
		planHandler,
		subscriptionHandler,
		serverHandler,
		nodeHandler,
		xrayHandler,
	)
	app.Router = appRouter
}
