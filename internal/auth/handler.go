package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service}
}

// LoginRequest represents the login request payload
// @Description Login request payload
type LoginRequest struct {
	Username string `json:"username" example:"admin"`                          // Username for authentication (optional if email provided)
	Email    string `json:"email" example:"admin@hotel.com"`                   // Email for authentication (optional if username provided)
	Password string `json:"password" binding:"required" example:"password123"` // Password for authentication
}

// LoginResponse represents the login response payload
// @Description Login response payload
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // JWT authentication token
}

// ErrorResponse represents an error response
// @Description Error response payload
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid credentials"` // Error message
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email or username and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Login credentials (email or username)"
// @Success 200 {object} LoginResponse "Login successful"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Determine which identifier to use (email or username)
	identifier := req.Email
	if identifier == "" {
		identifier = req.Username
	}

	// Validate that at least one identifier is provided
	if identifier == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Either email or username must be provided"})
		return
	}

	token, err := h.service.Login(c, identifier, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, ErrInvalidCredentials) || errors.Is(err, ErrUserInactive) {
			status = http.StatusUnauthorized
		}
		c.JSON(status, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: token})
}
