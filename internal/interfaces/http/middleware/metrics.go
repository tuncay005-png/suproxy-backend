package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/suproxy/backend/internal/infrastructure/metrics"
)

// MetricsMiddleware collects HTTP metrics for Prometheus
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Increment in-flight requests
		metrics.IncHTTPRequestsInFlight()
		defer metrics.DecHTTPRequestsInFlight()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Get path (use route pattern, not actual path to avoid high cardinality)
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		method := c.Request.Method
		status := strconv.Itoa(c.Writer.Status())

		// Record metrics
		metrics.RecordHTTPRequest(method, path, status)
		metrics.RecordHTTPDuration(method, path, duration)

		// Record errors (4xx and 5xx)
		if c.Writer.Status() >= 400 {
			metrics.RecordHTTPError(method, path, status)
		}
	}
}
