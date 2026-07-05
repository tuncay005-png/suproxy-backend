package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suproxy/backend/internal/infrastructure/logger"
	"github.com/suproxy/backend/internal/interfaces/http/response"
)

// ErrorHandler returns a middleware that handles panics and errors
func ErrorHandler(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("Panic recovered",
					"error", err,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
				)

				response.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
				c.Abort()
			}
		}()

		c.Next()

		// Check for errors after handlers
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			log.Error("Request error",
				"error", err.Error(),
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
			)
		}
	}
}
