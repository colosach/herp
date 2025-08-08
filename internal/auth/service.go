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

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if !user.IsActive {
		return "", ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", err
	}

	permissions, err := s.queries.GetUserPermissions(ctx, user.ID)
	if err != nil {
		return "", err
	}

	// permissionCodes := make([]string, len(permissions))
	// for i, p := range permissions {
	// 	permissionCodes[i] = p.Code
	// }
	token, err := jwt.GenerateToken(
		int(user.ID),
		user.Email,
		user.RoleName,
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
