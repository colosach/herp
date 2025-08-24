package auth

import (
	"context"
	"errors"
	"herp/internal/config"
	"herp/internal/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
	config *config.Config
}

func NewHandler(service *Service, c *config.Config) *Handler {
	return &Handler{service, c}
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
		log.Println("Error binding JSON:", err)
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
		log.Printf("Login error: %v", err)
		c.JSON(status, ErrorResponse{Error: err.Error()})
		if logErr := h.service.LogLogin(c.Request.Context(),
			req.Username,
			req.Email,
			c.ClientIP(),
			c.Request.UserAgent(),
			false,
			err.Error(),
		); logErr != nil {
			log.Printf("Failed to log login attempt: %v", logErr)
		}
		return
	}
	if logErr := h.service.LogLogin(context.Background(),
		req.Username,
		req.Email,
		c.ClientIP(),
		c.Request.UserAgent(),
		true,
		"",
	); logErr != nil {
		log.Printf("Failed to log login attempt: %v", logErr)
	}

	c.JSON(http.StatusOK, LoginResponse{Token: token})
}

// RegisterAdminRequest represents the login request payload
// @Description Register admin request payload
type RegisterAdminRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// Register godoc
// @Summary User Register
// @Description Create user with email, username, password and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body RegisterAdminRequest true "Register credentials (email, username and password)"
// @Success 200 {object} LoginResponse "Registration successful"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/auth/register [post]
func (h *Handler) RegisterAdmin(c *gin.Context) {
	var req RegisterAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "incorrect or empty register param"})
		return
	}
	// Generate verification code and expiry
	code := utils.GenerateOTP()
	expiry := time.Now().Add(10 * time.Minute)

	admin, err := h.service.RegisterAdmin(c, req.Username, req.Email, req.Password)
	if err != nil {
		log.Printf("error registering admin: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	err = h.service.SetEmailVerification(c.Request.Context(), admin.ID, code, expiry)
	if err != nil {
		log.Printf("error saving verification code: %v", err)
	}

	// Send verification email
	emailBody, _ := utils.RenderEmailTemplate("templates/auth/verify_email.html", map[string]any{
		"Username": admin.Username,
		"Code":     code,
	})
	plunk := utils.Plunk{HttpClient: http.DefaultClient, Config: h.config}
	err = plunk.SendEmail(admin.Email, "Verify your Herp account", emailBody)
	if err != nil {
		log.Printf("error sending verification email: %v", err)
	}

	c.JSON(http.StatusOK, admin)
}

// VerifyEmailRequest represents the login request payload
// @Description verify email request payload
type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

// Verify Email godoc
// @Summary Verify Admin Email
// @Description Verify admin email with email and code
// @Tags auth
// @Accept json
// @Produce json
// @Param body body VerifyEmailRequest true "Verify Email Request"
// @Success 200 {object} LoginResponse "Email verified"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/auth/verify-email [post]
func (h *Handler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ok, err := h.service.VerifyEmailCode(c.Request.Context(), req.Email, req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid or expired code"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Email verified successfully"})
}