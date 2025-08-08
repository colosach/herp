package auth

import (
	"database/sql"
	"errors"
	"fmt"
	db "herp/db/sqlc"
	"herp/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	service *Service
}

func NewAdminHandler(s *Service) *AdminHandler {
	return &AdminHandler{s}
}

func (h *AdminHandler) RegisterAdminRoutes(router *gin.RouterGroup, authSvc *Service) {
	admin := router.Group("/admin")
	admin.Use(AdminMiddleware(authSvc))

	// User management
	admin.GET("/users", h.ListUsers)
	admin.POST("/users", h.CreateUser)
	admin.GET("/users/:id", h.GetUser)
	admin.PUT("/users/:id", h.UpdateUser)
	admin.DELETE("/users/:id", h.DeleteUser)
	admin.POST("/users/:id/reset-password", h.ResetPassword)
	admin.GET("/users/:id/activity", h.GetUserActivityLogs)
	admin.GET("/users/:id/login-history", h.GetUserLoginHistory)

	// Role management
	admin.GET("/roles", h.ListRoles)
	admin.POST("/roles", h.CreateRole)
	admin.GET("/roles/:id", h.GetRole)
	admin.PUT("/roles/:id", h.UpdateRole)
	admin.DELETE("/roles/:id", h.DeleteRole)
	admin.POST("/roles/:id/permissions", h.AddPermissionToRole)
	admin.DELETE("/roles/:id/permissions/:permission_id", h.RemovePermissionFromRole)
	admin.GET("/roles/:id/permissions", h.GetRolePermissions)
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

// CreateUser creates a new user account
// @Summary Create a new user
// @Description Create a new user account with role assignment
// @Tags admin
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User creation data"
// @Success 201 {object} map[string]interface{} "User created successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/v1/admin/users [post]
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
		IsActive:     req.IsActive,
		Nin:          req.Nin,
		Gender:       req.Gender,
		DateOfBirth:  req.DateOfBirth,
	})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "user created successfully", user)
}

type UpdateUserRequest struct {
	FirstName   *string    `json:"first_name" binding:"required,min=2"`
	LastName    *string    `json:"last_name" binding:"required,min=2"`
	Email       *string    `json:"email"`
	RoleID      *int       `json:"role_id"`
	IsActive    *bool      `json:"is_active"`
	Nin         *string    `json:"nin" binding:"max=11"`
	Gender      *string    `json:"gender" binding:"required,oneof=male female"`
	DateOfBirth *time.Time `json:"date_of_birth" binding:"datetime=2006-01-02"`
}

// UpdateUser updates an existing user
// @Summary Update user information
// @Description Update user account information and settings
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body UpdateUserRequest true "User update data"
// @Success 200 {object} map[string]interface{} "User updated successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/v1/admin/users/{id} [put]
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var req UpdateUserRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	updateParams := db.UpdateUserParams{ID: int32(userID)}
	if req.FirstName != nil {
		updateParams.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		updateParams.LastName = *req.LastName
	}
	if req.RoleID != nil {
		updateParams.RoleID = int32(*req.RoleID)
	}
	if req.IsActive != nil {
		updateParams.IsActive = *req.IsActive
	}
	if req.Nin != nil {
		updateParams.Nin = *req.Nin
	}
	if req.DateOfBirth != nil {
		updateParams.DateOfBirth = *req.DateOfBirth
	}
	if req.Gender != nil {
		updateParams.Gender = *req.Gender
	}

	user, err := h.service.UpdateUser(c.Request.Context(), updateParams)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User data is updated", user)
}

// DeleteUser deletes a user account
// @Summary Delete user
// @Description Delete a user account from the system
// @Tags admin
// @Param id path int true "User ID"
// @Success 204 "User deleted successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/v1/admin/users/{id} [delete]
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.service.DeleteUser(c.Request.Context(), int32(userID)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusNoContent, "user is deleted", nil)
}

type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ResetPassword resets a user's password
// @Summary Reset user password
// @Description Reset password for a specific user
// @Tags admin
// @Accept json
// @Param id path int true "User ID"
// @Param password body ResetPasswordRequest true "New password data"
// @Success 204 "Password reset successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/v1/admin/users/{id}/reset-password [post]
func (h *AdminHandler) ResetPassword(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	params := db.UpdateUserPasswordParams{
		ID:           int32(userID),
		PasswordHash: req.NewPassword,
	}
	if err := h.service.ResetPassword(c.Request.Context(), params); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusNoContent, "password updated", nil)
}

// ListUsers retrieves all users
// @Summary List all users
// @Description Get a list of all users in the system
// @Tags admin
// @Produce json
// @Success 200 {array} map[string]interface{} "List of users"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/v1/admin/users [get]
func (h *AdminHandler) ListUsers(c *gin.Context) {
	users, err := h.service.queries.ListUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUser retrieves a specific user
// @Summary Get user by ID
// @Description Get detailed information about a specific user
// @Tags admin
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "User details"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/v1/admin/users/{id} [get]
func (h *AdminHandler) GetUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.service.queries.GetUserByID(c.Request.Context(), int32(userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ErrorResponse(c, http.StatusNotFound, "user not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", user)
}

// Role Management

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// CreateRole creates a new role
// @Summary Create a new role
// @Description Create a new role in the system
// @Tags admin
// @Accept json
// @Produce json
// @Param role body CreateRoleRequest true "Role creation data"
// @Success 201 {object} map[string]interface{} "Role created successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/v1/admin/roles [post]
func (h *AdminHandler) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	params := db.CreateRoleParams{
		Name:        req.Name,
		Description: sql.NullString{Valid: true, String: req.Description},
	}

	role, err := h.service.CreateRole(c.Request.Context(), params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "role created", role)
}

type UpdateRoleRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

// UpdateRole updates an existing role
// @Summary Update role
// @Description Update role information
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param role body UpdateRoleRequest true "Role update data"
// @Success 200 {object} map[string]interface{} "Role updated successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/v1/admin/roles/{id} [put]
func (h *AdminHandler) UpdateRole(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid role ID")
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateParams := db.UpdateRoleParams{ID: int32(roleID)}
	if req.Name != nil {
		updateParams.Name = *req.Name
	}
	if req.Description != nil {
		updateParams.Description = sql.NullString{Valid: true, String: *req.Description}
	}

	role, err := h.service.UpdateRole(c.Request.Context(), updateParams)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "role updated", role)
}

// DeleteRole deletes a role
// @Summary Delete role
// @Description Delete a role from the system
// @Tags admin
// @Param id path int true "Role ID"
// @Success 204 "Role deleted successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/v1/admin/roles/{id} [delete]
func (h *AdminHandler) DeleteRole(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid role id")
		return
	}

	if err := h.service.DeleteRole(c.Request.Context(), int32(roleID)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusNoContent, "role deleted", nil)
}

// ListRoles retrieves all roles
// @Summary List all roles
// @Description Get a list of all roles in the system
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]interface{} "List of roles"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/v1/admin/roles [get]
func (h *AdminHandler) ListRoles(c *gin.Context) {
	roles, err := h.service.queries.ListRoles(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", roles)
}

// GetRole retrieves a specific role
// @Summary Get role by ID
// @Description Get detailed information about a specific role
// @Tags admin
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} map[string]interface{} "Role details"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Role not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/v1/admin/roles/{id} [get]
func (h *AdminHandler) GetRole(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid role id")
		return
	}

	role, err := h.service.queries.GetRoleByID(c.Request.Context(), int32(roleID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ErrorResponse(c, http.StatusNotFound, "role not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", role)
}

type ManageRolePermissionRequest struct {
	PermissionID int `json:"permission_id" binding:"required"`
}

func (h *AdminHandler) AddPermissionToRole(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid role ID")
		return
	}

	var req ManageRolePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	params := db.AddPermissionToRoleParams{
		RoleID:       int32(roleID),
		PermissionID: int32(req.PermissionID),
	}

	if err := h.service.AddPermissionToRole(c.Request.Context(), params); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusNoContent, fmt.Sprintf("permission %d added to role %d", roleID, req.PermissionID), nil)
}

func (h *AdminHandler) RemovePermissionFromRole(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid role ID")
		return
	}

	permissionID, err := strconv.Atoi(c.Param("permission_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid permission ID")
		return
	}

	params := db.RemovePermissionFromRoleParams{
		RoleID:       int32(roleID),
		PermissionID: int32(permissionID),
	}

	if err := h.service.RemovePermissionFromRole(c.Request.Context(), params); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusNoContent, fmt.Sprintf("permission %d removed from role %d", permissionID, roleID), nil)
}

func (h *AdminHandler) GetRolePermissions(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid role ID")
		return
	}

	permissions, err := h.service.queries.GetRolePermissions(c.Request.Context(), int32(roleID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", permissions)
}

// Activity Logs

func (h *AdminHandler) GetUserActivityLogs(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if limit > 1000 {
		limit = 1000
	}

	logs, err := h.service.queries.GetUserActivityLogs(c.Request.Context(), db.GetUserActivityLogsParams{
		UserID: int32(userID),
		Limit:  int32(limit),
	})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", logs)
}

func (h *AdminHandler) GetUserLoginHistory(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if limit > 1000 {
		limit = 1000
	}

	params := db.GetUserLoginHistoryParams{
		UserID: int32(userID),
		Limit:  int32(limit),
	}

	history, err := h.service.queries.GetUserLoginHistory(c.Request.Context(), params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", history)
}
