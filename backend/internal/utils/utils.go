package utils

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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

// UploadFile validates and saves an uploaded file.
// Returns the relative URL path (e.g. /images/123_logo.png) or an error.
func UploadFile(c *gin.Context, fieldName string, saveDir string, maxSize int64) (string, error) {
	file, err := c.FormFile(fieldName)
	if err != nil {
		// No file provided
		return "", err
	}

	// Check file size
	if file.Size > maxSize {
		return "", fmt.Errorf("file too large, max %d bytes allowed", maxSize)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExt := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	if !allowedExt[ext] {
		return "", fmt.Errorf("invalid file extension: only JPG/PNG allowed")
	}

	// Check MIME type
	openedFile, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("could not open uploaded file: %v", err)
	}
	defer openedFile.Close()

	buffer := make([]byte, 512)
	if _, err := openedFile.Read(buffer); err != nil {
		return "", fmt.Errorf("could not read uploaded file: %v", err)
	}

	contentType := http.DetectContentType(buffer)
	allowedMime := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}
	if !allowedMime[contentType] {
		return "", fmt.Errorf("invalid file type: only JPG/PNG allowed")
	}

	// Ensure save directory exists
	if _, statErr := os.Stat(saveDir); os.IsNotExist(statErr) {
		os.MkdirAll(saveDir, os.ModePerm)
	}

	// Generate unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
	filePath := filepath.Join(saveDir, filename)

	// Save file
	if saveErr := c.SaveUploadedFile(file, filePath); saveErr != nil {
		return "", fmt.Errorf("could not save file: %v", saveErr)
	}

	// Return relative URL for serving via Gin Static
	return "/" + filePath, nil
}
