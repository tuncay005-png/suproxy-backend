package bootstrap

import (
	"github.com/gin-gonic/gin"
	authUseCase "github.com/suproxy/backend/internal/application/usecase/auth"
	"github.com/suproxy/backend/internal/interfaces/http/handler"
	"github.com/suproxy/backend/internal/interfaces/http/router"
	"github.com/suproxy/backend/internal/infrastructure/repository"
)

func InitializeAuthSystem(app *Application, engine *gin.Engine) {
	// Initialize repositories
	userRepo := repository.NewUserRepository(app.Database.DB)

	// Initialize use cases
	registerCmd := authUseCase.NewRegisterCommand(userRepo, app.Logger)
	loginCmd := authUseCase.NewLoginCommand(userRepo, app.JWTManager, app.Logger)
	refreshTokenCmd := authUseCase.NewRefreshTokenCommand(userRepo, app.JWTManager, app.Logger)
	getCurrentUserQuery := authUseCase.NewGetCurrentUserQuery(userRepo, app.Logger)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(registerCmd, loginCmd, refreshTokenCmd, getCurrentUserQuery)

	// Initialize router
	appRouter := router.NewRouter(engine, authHandler, app.JWTManager)
	app.Router = appRouter
}
