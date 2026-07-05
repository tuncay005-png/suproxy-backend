package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

// RequestLogger returns a middleware that logs HTTP requests
func RequestLogger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		userAgent := c.Request.UserAgent()

		log.Info("HTTP Request",
			"method", method,
			"path", path,
			"query", query,
			"status", statusCode,
			"latency", latency,
			"ip", clientIP,
			"user_agent", userAgent,
		)
	}
}
