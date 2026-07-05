package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
	})
}

func SuccessCreated(c *gin.Context, data interface{}) {
	Success(c, http.StatusCreated, data)
}

func SuccessOK(c *gin.Context, data interface{}) {
	Success(c, http.StatusOK, data)
}

func SuccessNoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, Response{
		Success: true,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error: &Error{
			Code:    code,
			Message: message,
		},
	})
}

func BadRequest(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", message)
}

func Unauthorized(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

func Forbidden(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, "FORBIDDEN", message)
}

func NotFound(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", message)
}

func Conflict(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusConflict, "CONFLICT", message)
}

func InternalError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", message)
}

func ValidationError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", message)
}
