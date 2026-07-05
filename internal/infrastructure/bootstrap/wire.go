package bootstrap

import (
	"github.com/gin-gonic/gin"
	authUseCase "github.com/suproxy/backend/internal/application/usecase/auth"
	"github.com/suproxy/backend/internal/infrastructure/repository"
	"github.com/suproxy/backend/internal/interfaces/http/handler"
	"github.com/suproxy/backend/internal/interfaces/http/router"
)

func InitializeAuthSystem(app *Application, engine *gin.Engine) {
	// Initialize repositories
	userRepo := repository.NewUserRepository(app.Database.DB)
	refreshTokenRepo := repository.NewRefreshTokenRepository(app.Database.DB)
	auditRepo := repository.NewAuditLogRepository(app.Database.DB)

	// Initialize use cases
	registerCmd := authUseCase.NewRegisterCommand(userRepo, app.Logger)
	loginCmd := authUseCase.NewLoginCommand(userRepo, refreshTokenRepo, auditRepo, app.JWTManager, app.Logger)
	refreshTokenCmd := authUseCase.NewRefreshTokenCommand(userRepo, refreshTokenRepo, auditRepo, app.JWTManager, app.Logger)
	logoutCmd := authUseCase.NewLogoutCommand(refreshTokenRepo, auditRepo, app.Logger)
	getCurrentUserQuery := authUseCase.NewGetCurrentUserQuery(userRepo, app.Logger)
	getSessionsQuery := authUseCase.NewGetSessionsQuery(refreshTokenRepo, app.Logger)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(registerCmd, loginCmd, refreshTokenCmd, logoutCmd, getCurrentUserQuery, getSessionsQuery)

	// Initialize router
	appRouter := router.NewRouter(engine, authHandler, app.JWTManager)
	app.Router = appRouter
}
