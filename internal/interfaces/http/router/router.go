package router

import (
	"github.com/gin-gonic/gin"
	"github.com/suproxy/backend/internal/infrastructure/jwt"
	"github.com/suproxy/backend/internal/infrastructure/logger"
	"github.com/suproxy/backend/internal/interfaces/http/handler"
	"github.com/suproxy/backend/internal/interfaces/http/middleware"
)

type Router struct {
	engine                 *gin.Engine
	logger                 *logger.Logger
	jwtManager             *jwt.Manager
	healthHandler          *handler.HealthHandler
	authHandler            *handler.AuthHandler
	userHandler            *handler.UserHandler
	planHandler            *handler.PlanHandler
	subscriptionHandler    *handler.SubscriptionHandler
	serverHandler          *handler.ServerHandler
	nodeHandler            *handler.NodeHandler
	xrayHandler            *handler.XrayHandler
}

func NewRouter(
	engine *gin.Engine,
	log *logger.Logger,
	jwtManager *jwt.Manager,
	healthHandler *handler.HealthHandler,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	planHandler *handler.PlanHandler,
	subscriptionHandler *handler.SubscriptionHandler,
	serverHandler *handler.ServerHandler,
	nodeHandler *handler.NodeHandler,
	xrayHandler *handler.XrayHandler,
) *Router {
	return &Router{
		engine:              engine,
		logger:              log,
		jwtManager:          jwtManager,
		healthHandler:       healthHandler,
		authHandler:         authHandler,
		userHandler:         userHandler,
		planHandler:         planHandler,
		subscriptionHandler: subscriptionHandler,
		serverHandler:       serverHandler,
		nodeHandler:         nodeHandler,
		xrayHandler:         xrayHandler,
	}
}

func (r *Router) Setup() {
	// Global middlewares
	r.engine.Use(middleware.CORS())
	r.engine.Use(middleware.ErrorHandler(r.logger))
	r.engine.Use(middleware.RequestLogger(r.logger))

	// Health check endpoints (no auth required)
	r.engine.GET("/health", r.healthHandler.Health)
	r.engine.GET("/ready", r.healthHandler.Ready)

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

			// Protected auth routes
			authenticated := auth.Group("")
			authenticated.Use(middleware.AuthMiddleware(r.jwtManager))
			{
				authenticated.GET("/me", r.authHandler.GetCurrentUser)
				authenticated.GET("/sessions", r.authHandler.GetSessions)
				authenticated.DELETE("/sessions/:id", r.authHandler.LogoutSingle)
				authenticated.POST("/logout-all", r.authHandler.LogoutAll)
			}
		}

		// User routes (protected)
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			users.GET("/me", r.userHandler.GetMe)
			users.PUT("/me", r.userHandler.UpdateMe)
		}

		// Plan routes (public)
		plans := v1.Group("/plans")
		{
			plans.GET("", r.planHandler.ListPlans)
			plans.GET("/:id", r.planHandler.GetPlan)
		}

		// Subscription routes (protected)
		subscriptions := v1.Group("/subscriptions")
		subscriptions.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			subscriptions.GET("/me", r.subscriptionHandler.GetMySubscription)
		}

		// Server routes (protected)
		servers := v1.Group("/servers")
		servers.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			servers.GET("", r.serverHandler.ListServers)
		}

		// Node routes (protected)
		nodes := v1.Group("/nodes")
		nodes.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			nodes.GET("", r.nodeHandler.ListNodes)
		}

		// Xray routes (protected)
		xray := v1.Group("/xray")
		xray.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			xray.GET("/instances", r.xrayHandler.ListInstances)
		}
	}
}
