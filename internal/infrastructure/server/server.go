package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/suproxy/backend/internal/infrastructure/config"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type Server struct {
	config *config.Config
	logger *logger.Logger
	router *gin.Engine
	server *http.Server
}

func New(cfg *config.Config, log *logger.Logger) *Server {
	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(loggerMiddleware(log))

	s := &Server{
		config: cfg,
		logger: log,
		router: router,
	}

	s.setupRoutes()

	s.server = &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	return s
}

func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.GET("/health", s.healthCheckHandler)

	// API v1 group
	v1 := s.router.Group("/api/v1")
	{
		v1.GET("/ping", s.pingHandler)
	}
}

func (s *Server) healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "suproxy-backend",
		"version": "1.0.0",
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

func loggerMiddleware(log *logger.Logger) gin.HandlerFunc {
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

		log.Infow("HTTP Request",
			"status", statusCode,
			"method", method,
			"path", path,
			"ip", clientIP,
			"latency", fmt.Sprintf("%v", latency),
			"error", errorMessage,
		)
	}
}
