package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/suproxy/backend/internal/infrastructure/bootstrap"
)

type Server struct {
	app    *bootstrap.Application
	router *gin.Engine
	server *http.Server
}

func New(app *bootstrap.Application) *Server {
	// Set Gin mode
	if app.Config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(loggerMiddleware(app))

	// Initialize auth system
	bootstrap.InitializeAuthSystem(app, router)

	s := &Server{
		app:    app,
		router: router,
	}

	s.setupRoutes()

	s.server = &http.Server{
		Addr:         app.Config.Server.Address,
		Handler:      router,
		ReadTimeout:  time.Duration(app.Config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(app.Config.Server.WriteTimeout) * time.Second,
	}

	return s
}

func (s *Server) setupRoutes() {
	// Health check endpoints
	s.router.GET("/health", s.healthCheckHandler)
	s.router.GET("/health/db", s.dbHealthCheckHandler)

	// API v1 group
	v1 := s.router.Group("/api/v1")
	{
		v1.GET("/ping", s.pingHandler)
	}
	
	// Setup additional routes via router
	if s.app.Router != nil {
		s.app.Router.Setup()
	}
}

func (s *Server) healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "suproxy-backend",
		"version": "1.0.0",
	})
}

func (s *Server) dbHealthCheckHandler(c *gin.Context) {
	if err := s.app.Database.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  "database connection failed",
		})
		return
	}

	stats := s.app.Database.GetStats()
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"stats":  stats,
	})
}

func (s *Server) pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func loggerMiddleware(app *bootstrap.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		app.Logger.Infow("HTTP Request",
			"status", statusCode,
			"method", method,
			"path", path,
			"ip", clientIP,
			"latency", fmt.Sprintf("%v", latency),
			"error", errorMessage,
		)
	}
}
