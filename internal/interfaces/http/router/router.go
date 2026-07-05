package router

import (
	"github.com/gin-gonic/gin"
	"github.com/suproxy/backend/internal/interfaces/http/handler"
	"github.com/suproxy/backend/internal/interfaces/http/middleware"
	"github.com/suproxy/backend/internal/infrastructure/jwt"
)

type Router struct {
	engine      *gin.Engine
	authHandler *handler.AuthHandler
	jwtManager  *jwt.Manager
}

func NewRouter(engine *gin.Engine, authHandler *handler.AuthHandler, jwtManager *jwt.Manager) *Router {
	return &Router{
		engine:      engine,
		authHandler: authHandler,
		jwtManager:  jwtManager,
	}
}

func (r *Router) Setup() {
	// API v1 group
	v1 := r.engine.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", r.authHandler.Register)
			auth.POST("/login", r.authHandler.Login)
			auth.POST("/refresh", r.authHandler.RefreshToken)
			auth.POST("/logout", r.authHandler.Logout)
			
			// Protected routes
			authenticated := auth.Group("")
			authenticated.Use(middleware.AuthMiddleware(r.jwtManager))
			{
				authenticated.GET("/me", r.authHandler.GetCurrentUser)
			}
		}
	}
}
