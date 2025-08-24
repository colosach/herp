package middleware

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// NewRequestLogger returns a Gin middleware that logs request/response details
// to both stdout and the specified file. The directory for the log file will be
// created if it does not exist.
func NewRequestLogger(logFilePath string) gin.HandlerFunc {
	// Ensure the directory exists
	dir := filepath.Dir(logFilePath)
	if dir != "." && dir != "" {
		_ = os.MkdirAll(dir, 0o755)
	}

	// Open or create the log file for appending
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		// Fallback to stdout-only if file cannot be opened
		return func(c *gin.Context) {
			start := time.Now()
			c.Next()
			writeLog(os.Stdout, c, start)
		}
	}

	multiDest := io.MultiWriter(os.Stdout, file)

	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		writeLog(multiDest, c, start)
	}
}

func writeLog(w io.Writer, c *gin.Context, start time.Time) {
	latency := time.Since(start)
	statusCode := c.Writer.Status()
	clientIP := c.ClientIP()
	method := c.Request.Method
	path := c.FullPath()
	if path == "" {
		path = c.Request.URL.Path
	}
	userAgent := c.Request.UserAgent()
	bytesWritten := c.Writer.Size()
	if bytesWritten < 0 {
		bytesWritten = 0
	}

	// Capture error information
	var errorDetails []string

	// Capture Gin validation errors
	if len(c.Errors) > 0 {
		for _, err := range c.Errors {
			errorDetails = append(errorDetails, err.Error())
		}
	}

	// Capture HTTP status errors (4xx, 5xx)
	if statusCode >= 400 {
		// Try to capture response body for error details
		if c.Writer.Size() > 0 {
			// Note: We can't easily capture response body in middleware
			// but we can log the status code and any error messages
			errorDetails = append(errorDetails, fmt.Sprintf("HTTP %d", statusCode))
		}
	}

	// Capture specific error types based on status codes
	switch {
	case statusCode >= 500:
		errorDetails = append(errorDetails, "server_error")
	case statusCode == 404:
		errorDetails = append(errorDetails, "not_found")
	case statusCode == 401:
		errorDetails = append(errorDetails, "unauthorized")
	case statusCode == 403:
		errorDetails = append(errorDetails, "forbidden")
	case statusCode == 400:
		errorDetails = append(errorDetails, "bad_request")
	case statusCode == 422:
		errorDetails = append(errorDetails, "validation_error")
	}

	// Common log format: time, status, method, path, latency, size, ip, ua, error
	timestamp := time.Now().Format(time.RFC3339)
	_, _ = fmt.Fprintf(
		w,
		"%s | %3d | %13v | %7s | %-40s | %6dB | ip=%s | ua=\"%s\"%s\n",
		timestamp,
		statusCode,
		latency,
		method,
		path,
		bytesWritten,
		clientIP,
		userAgent,
		formatErrors(errorDetails),
	)

	// Log request ID if present
	if reqID := c.GetHeader("X-Request-ID"); reqID != "" {
		_, _ = fmt.Fprintf(w, "\trequest_id=%s\n", reqID)
	}

	// Log additional error details if any
	if len(errorDetails) > 0 {
		_, _ = fmt.Fprintf(w, "\terror_details=%s\n", strings.Join(errorDetails, ", "))
	}
}

func formatErrors(errors []string) string {
	if len(errors) == 0 {
		return ""
	}
	return " | error=\"" + strings.Join(errors, "; ") + "\""
}
