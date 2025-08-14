package auth

import (
	"context"
	"database/sql"
	"errors"
	db "herp/db/sqlc"
	"herp/pkg/jwt"
	"slices"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserInactive       = errors.New("user is inactive")
)

type Service struct {
	queries   *db.Queries
	jwtSecret string
	jwtExpiry time.Duration
}

func NewService(queries *db.Queries, jwtSecret string, jwtExpiry time.Duration) *Service {
	return &Service{
		queries:   queries,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

func (s *Service) Login(ctx context.Context, emailOrUsername, password string) (string, error) {
	// Try to find user by email first, then by username
	userByEmail, err := s.queries.GetUserByEmail(ctx, emailOrUsername)
	if err != nil {
		// If not found by email, try by username
		userByUsername, err := s.queries.GetUserByUsername(ctx, emailOrUsername)
		if err != nil {
			return "", ErrInvalidCredentials
		}
		// Convert to common struct for consistent handling
		if !userByUsername.IsActive {
			return "", ErrUserInactive
		}

		if err := bcrypt.CompareHashAndPassword([]byte(userByUsername.PasswordHash), []byte(password)); err != nil {
			return "", ErrInvalidCredentials
		}

		permissions, err := s.queries.GetUserPermissions(ctx, userByUsername.ID)
		if err != nil {
			return "", err
		}

		token, err := jwt.GenerateToken(
			int(userByUsername.ID),
			userByUsername.Email,
			userByUsername.RoleName,
			s.jwtSecret,
			permissions,
			s.jwtExpiry,
		)
		if err != nil {
			return "", err
		}

		return token, nil
	}

	// User found by email
	if !userByEmail.IsActive {
		return "", ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userByEmail.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	permissions, err := s.queries.GetUserPermissions(ctx, userByEmail.ID)
	if err != nil {
		return "", err
	}

	token, err := jwt.GenerateToken(
		int(userByEmail.ID),
		userByEmail.Email,
		userByEmail.RoleName,
		s.jwtSecret,
		permissions,
		s.jwtExpiry,
	)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) HasPermission(claims *jwt.Claims, requiredPermission string) bool {
	return slices.Contains(claims.Permissions, requiredPermission)
}

// Admin user management functions
func (s *Service) CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return db.User{}, err
	}

	params.PasswordHash = string(hashedPassword)
	return s.queries.CreateUser(ctx, params)
}

func (s *Service) UpdateUser(ctx context.Context, params db.UpdateUserParams) (db.User, error) {
	return s.queries.UpdateUser(ctx, params)
}

func (s *Service) DeleteUser(ctx context.Context, id int32) error {
	return s.queries.DeleteUser(ctx, id)
}

func (s *Service) ResetPassword(ctx context.Context, params db.UpdateUserPasswordParams) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	params.PasswordHash = string(hashedPassword)
	return s.queries.UpdateUserPassword(ctx, params)
}

// Role management functions
func (s *Service) CreateRole(ctx context.Context, params db.CreateRoleParams) (db.Role, error) {
	return s.queries.CreateRole(ctx, params)
}

func (s *Service) UpdateRole(ctx context.Context, params db.UpdateRoleParams) (db.Role, error) {
	return s.queries.UpdateRole(ctx, params)
}

func (s *Service) DeleteRole(ctx context.Context, id int32) error {
	return s.queries.DeleteRole(ctx, id)
}

func (s *Service) AddPermissionToRole(ctx context.Context, params db.AddPermissionToRoleParams) error {
	return s.queries.AddPermissionToRole(ctx, params)
}

func (s *Service) RemovePermissionFromRole(ctx context.Context, params db.RemovePermissionFromRoleParams) error {
	return s.queries.RemovePermissionFromRole(ctx, params)
}

// Logging functions
// Additional helper functions
func (s *Service) GetUserByID(ctx context.Context, id int32) (db.GetUserByIDRow, error) {
	return s.queries.GetUserByID(ctx, id)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (db.GetUserByEmailRow, error) {
	return s.queries.GetUserByEmail(ctx, email)
}

func (s *Service) GetUserByUsername(ctx context.Context, username string) (db.GetUserByUsernameRow, error) {
	return s.queries.GetUserByUsername(ctx, username)
}

func (s *Service) ListUsers(ctx context.Context) ([]db.ListUsersRow, error) {
	return s.queries.ListUsers(ctx)
}

func (s *Service) ListRoles(ctx context.Context) ([]db.Role, error) {
	return s.queries.ListRoles(ctx)
}

func (s *Service) GetRolePermissions(ctx context.Context, roleID int32) ([]db.Permission, error) {
	return s.queries.GetRolePermissions(ctx, roleID)
}

// Logging functions
func (s *Service) LogUserActivity(ctx context.Context, userID int, action, description, ip, userAgent string) error {
	_, err := s.queries.LogUserActivity(ctx, db.LogUserActivityParams{
		UserID:      int32(userID),
		Action:      action,
		Description: description,
		IpAddress:   sql.NullString{Valid: true, String: ip},
		UserAgent:   sql.NullString{Valid: true, String: userAgent},
	})
	return err
}

func (s *Service) LogLoginAttempt(ctx context.Context, userID int, ip, userAgent string, success bool) error {
	_, err := s.queries.LogLoginAttempt(ctx, db.LogLoginAttemptParams{
		UserID:    int32(userID),
		IpAddress: sql.NullString{Valid: true, String: ip},
		UserAgent: sql.NullString{Valid: true, String: userAgent},
		Success:   success,
	})
	return err
}
