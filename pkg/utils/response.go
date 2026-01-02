package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Response represents the standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	TraceID string      `json:"trace_id"`
}

// SuccessResponse sends a successful response
func SuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	traceID := getTraceID(c)
	
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
		TraceID: traceID,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, err error) {
	traceID := getTraceID(c)
	
	c.JSON(statusCode, Response{
		Success: false,
		Error:   err.Error(),
		TraceID: traceID,
	})
}

// ErrorMessageResponse sends an error response with a custom message
func ErrorMessageResponse(c *gin.Context, statusCode int, message string) {
	traceID := getTraceID(c)
	
	c.JSON(statusCode, Response{
		Success: false,
		Error:   message,
		TraceID: traceID,
	})
}

// getTraceID retrieves or generates a trace ID for the request
func getTraceID(c *gin.Context) string {
	// Try to get from context first (set by middleware)
	if traceID, exists := c.Get("trace_id"); exists {
		return traceID.(string)
	}
	
	// Generate new trace ID if not exists
	return uuid.New().String()
}
