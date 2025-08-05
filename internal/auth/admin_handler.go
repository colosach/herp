package auth

import (
	"database/sql"
	db "herp/db/sqlc"
	"herp/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	service *Service
}

func NewAdminHandler(s *Service) *AdminHandler {
	return &AdminHandler{s}
}

// User Management
type CreateUserRequest struct {
	FirstName   string    `json:"first_name" binding:"required,min=2"`
	LastName    string    `json:"last_name" binding:"required,min=2"`
	Email       string    `json:"email" binding:"required,email"`
	Password    string    `json:"password" binding:"required,min=4"`
	RoleID      int       `json:"role_id" binding:"required"`
	IsActive    bool      `json:"is_active" binding:"required"`
	Nin         string    `json:"nin" binding:"max=11"`
	Gender      string    `json:"gender" binding:"required,oneof=male female"`
	DateOfBirth time.Time `json:"date_of_birth" binding:"datetime=2006-01-02"`
}

func (h *AdminHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.service.CreateUser(c, db.CreateUserParams{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PasswordHash: req.Password,
		RoleID:       int32(req.RoleID),
		IsActive:     sql.NullBool{Valid: true, Bool: req.IsActive},
		Nin:          sql.NullString{String: req.Nin, Valid: true},
		Gender:       sql.NullString{String: req.Gender, Valid: true},
		DateOfBirth:  sql.NullTime{Time: req.DateOfBirth, Valid: true},
	})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "user created successfully", user)
}
