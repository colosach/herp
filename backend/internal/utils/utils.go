package utils

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Get client IP from context
func GetClientIP(c context.Context) string {
	if ginCtx, ok := c.(*gin.Context); ok {
		// Try to get IP from X-Forwarded-For header (if behind proxy)
		if ip := ginCtx.GetHeader("X-Forwarded-For"); ip != "" {
			return strings.Split(ip, ",")[0] // First IP in chain
		}
		// Fall back to remote address
		return ginCtx.ClientIP()
	}
	return "unknown"
}

func WriteActivityDetails(username, email, action string, time time.Time) string {
	return fmt.Sprintf("User %s with email %s performed action: %s at %s", username, email, action, time)
}