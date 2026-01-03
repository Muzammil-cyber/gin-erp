package middleware

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ANSI color codes
const (
	reset   = "\033[0m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[37m"
	gray    = "\033[90m"
)

var (
	logFile       *os.File
	consoleLogger *log.Logger
	fileLogger    *log.Logger
)

// InitLogger initializes the log file
func InitLogger() error {
	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil { // nolint:gosec
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Create log file with date
	logFileName := fmt.Sprintf("app-%s.log", time.Now().Format("2006-01-02"))
	logFilePath := filepath.Join(logsDir, logFileName)

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666) // nolint:gosec
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	logFile = file

	// Create separate loggers
	consoleLogger = log.New(os.Stdout, "", 0)
	fileLogger = log.New(logFile, "", 0)

	return nil
}

// CloseLogger closes the log file
func CloseLogger() {
	if logFile != nil {
		logFile.Close() // nolint:errcheck,gosec
	}
}

// colorizeMethod returns colored method string
func colorizeMethod(method string) string {
	switch method {
	case "GET":
		return green + method + reset
	case "POST":
		return blue + method + reset
	case "PUT":
		return yellow + method + reset
	case "DELETE":
		return red + method + reset
	case "PATCH":
		return magenta + method + reset
	default:
		return cyan + method + reset
	}
}

// colorizeStatus returns colored status code
func colorizeStatus(status int) string {
	switch {
	case status >= 200 && status < 300:
		return green
	case status >= 300 && status < 400:
		return yellow
	case status >= 400 && status < 500:
		return red
	case status >= 500:
		return red
	default:
		return white
	}
}

// responseWriter wraps gin.ResponseWriter to capture response body
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// TraceIDMiddleware adds a trace ID to each request
func TraceIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := uuid.New().String()
		c.Set("trace_id", traceID)
		c.Header("X-Trace-ID", traceID)
		c.Next()
	}
}

// LoggerMiddleware logs request and response details
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		traceID, _ := c.Get("trace_id")

		// Read request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// Restore the body for handlers
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Wrap response writer to capture response body
		blw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		c.Next()

		// Calculate duration
		duration := time.Since(start)
		statusCode := c.Writer.Status()
		timestamp := start.Format("2006-01-02 15:04:05")

		// Console log - simple and colorized
		consoleLogger.Printf("%s [%s] %s %s | Status: %s%d%s | Duration: %v | TraceID: %s",
			timestamp,
			colorizeMethod(method),
			path,
			c.ClientIP(),
			colorizeStatus(statusCode),
			statusCode,
			reset,
			duration,
			traceID,
		)

		// File log - detailed without colors
		fileLogger.Printf("\n═══════════════════════════════════════════════════════════════")
		fileLogger.Printf("▶ REQUEST | TraceID: %s", traceID)
		fileLogger.Printf("  %s %s", method, path)
		fileLogger.Printf("  Client: %s", c.ClientIP())
		fileLogger.Printf("  Time: %s", timestamp)

		// Log important headers
		if auth := c.GetHeader("Authorization"); auth != "" {
			fileLogger.Printf("  Auth: Bearer ***")
		}
		if contentType := c.GetHeader("Content-Type"); contentType != "" {
			fileLogger.Printf("  Content-Type: %s", contentType)
		}

		// Log request body if exists (limit to 500 chars)
		if len(requestBody) > 0 {
			bodyStr := string(requestBody)
			if len(bodyStr) > 500 {
				bodyStr = bodyStr[:500] + "..."
			}
			fileLogger.Printf("  Body: %s", bodyStr)
		}

		// Log response
		fileLogger.Printf("◀ RESPONSE | TraceID: %s", traceID)
		fileLogger.Printf("  Status: %d", statusCode)
		fileLogger.Printf("  Duration: %v", duration)
		fileLogger.Printf("  Size: %d bytes", blw.body.Len())

		// Log response body if exists (limit to 500 chars)
		if blw.body.Len() > 0 {
			bodyStr := blw.body.String()
			if len(bodyStr) > 500 {
				bodyStr = bodyStr[:500] + "..."
			}
			fileLogger.Printf("  Body: %s", bodyStr)
		}

		// Log errors if any
		if len(c.Errors) > 0 {
			fileLogger.Printf("  Errors: %v", c.Errors)
		}

		fileLogger.Printf("═══════════════════════════════════════════════════════════════\n")
	}
}
