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
	adminHandler           *handler.AdminHandler
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
	adminHandler *handler.AdminHandler,
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
		adminHandler:        adminHandler,
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

		// Plan routes (public read, admin write)
		plans := v1.Group("/plans")
		{
			// Public endpoints
			plans.GET("", r.planHandler.ListPlans)
			plans.GET("/:id", r.planHandler.GetPlan)

			// Admin endpoints (future implementation)
			// admin := plans.Group("")
			// admin.Use(middleware.AuthMiddleware(r.jwtManager))
			// admin.Use(middleware.RequireAdmin())
			// {
			// 	admin.POST("", r.planHandler.CreatePlan)
			// 	admin.PUT("/:id", r.planHandler.UpdatePlan)
			// 	admin.DELETE("/:id", r.planHandler.DeletePlan)
			// }
		}

		// Subscription routes (protected - user)
		subscriptions := v1.Group("/subscriptions")
		subscriptions.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			subscriptions.GET("/me", r.subscriptionHandler.GetMySubscription)
		}

		// Server routes (protected - user read, admin write)
		servers := v1.Group("/servers")
		servers.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			servers.GET("", r.serverHandler.ListServers)
		}

		// Node routes (protected - user read, admin write)
		nodes := v1.Group("/nodes")
		nodes.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			nodes.GET("", r.nodeHandler.ListNodes)
		}

		// Xray routes (protected - user read, admin write)
		xray := v1.Group("/xray")
		xray.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			xray.GET("/instances", r.xrayHandler.ListInstances)
		}

		// Admin routes (protected - admin only)
		// Foundation for all admin operations
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(r.jwtManager))
		admin.Use(middleware.AdminAuthorization(r.logger))
		{
			// Health check endpoint - verifies admin API is accessible
			admin.GET("/health", r.adminHandler.HealthCheck)

			// User Management (Phase 17.2)
			users := admin.Group("/users")
			{
				users.GET("", r.adminHandler.ListUsers)           // List users with filters
				users.GET("/:id", r.adminHandler.GetUser)         // Get user details
				users.PUT("/:id/status", r.adminHandler.UpdateUserStatus) // Update user status
				users.PUT("/:id/role", r.adminHandler.UpdateUserRole)     // Update user role
			}

			// Xray Instance Management (Phase 17.3)
			xrayGroup := admin.Group("/xray")
			{
				instances := xrayGroup.Group("/instances")
				{
					instances.GET("", r.adminHandler.ListInstances)
					instances.GET("/:id", r.adminHandler.GetInstance)
					instances.POST("/:id/start", r.adminHandler.StartInstance)
					instances.POST("/:id/stop", r.adminHandler.StopInstance)
					instances.POST("/:id/restart", r.adminHandler.RestartInstance)
					instances.POST("/:id/reload", r.adminHandler.ReloadInstance)
					instances.GET("/:id/health", r.adminHandler.CheckInstanceHealth)
					instances.GET("/:id/stats", r.adminHandler.GetInstanceStats)
				}

				// Inbound Management (Phase 17.4)
				inbounds := xrayGroup.Group("/inbounds")
				{
					inbounds.GET("", r.adminHandler.ListInbounds)
					inbounds.GET("/:id", r.adminHandler.GetInbound)
					inbounds.POST("", r.adminHandler.CreateInbound)
					inbounds.PUT("/:id", r.adminHandler.UpdateInbound)
					inbounds.DELETE("/:id", r.adminHandler.DeleteInbound)
					inbounds.PUT("/:id/enable", r.adminHandler.EnableInbound)
					inbounds.PUT("/:id/disable", r.adminHandler.DisableInbound)
				}

				// Client Management (Phase 17.5)
				clients := xrayGroup.Group("/clients")
				{
					clients.GET("", r.adminHandler.ListClients)
					clients.GET("/:id", r.adminHandler.GetClient)
					clients.POST("", r.adminHandler.CreateClient)
					clients.DELETE("/:id", r.adminHandler.DeleteClient)
					clients.PUT("/:id/enable", r.adminHandler.EnableClient)
					clients.PUT("/:id/disable", r.adminHandler.DisableClient)
					clients.POST("/:id/regenerate-uuid", r.adminHandler.RegenerateClientUUID)
					clients.POST("/:id/reprovision", r.adminHandler.ReprovisionClient)
				}
			}

			// Audit Log Management (Phase 17.6)
			audit := admin.Group("/audit")
			{
				audit.GET("/logs", r.adminHandler.ListAuditLogs)
				audit.GET("/logs/:id", r.adminHandler.GetAuditLog)
				audit.GET("/stats", r.adminHandler.GetAuditStats)
			}

			// System Admin (Phase 17.7)
			system := admin.Group("/system")
			{
				system.GET("/health", r.adminHandler.GetSystemHealth)
				system.GET("/stats", r.adminHandler.GetSystemStats)
				system.GET("/version", r.adminHandler.GetVersion)
				system.GET("/database", r.adminHandler.GetDatabaseStatus)
				system.GET("/xray", r.adminHandler.GetXraySystemStatus)
			}
		}
	}
}
