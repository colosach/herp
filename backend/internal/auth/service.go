package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	db "herp/db/sqlc"
	"herp/pkg/jwt"
	"herp/pkg/redis"
	"log"
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
	redis     *redis.Redis
}

func NewService(queries *db.Queries, jwtSecret string, jwtExpiry time.Duration, redis *redis.Redis) *Service {
	return &Service{
		queries:   queries,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
		redis:     redis,
	}
}

func (s *Service) RegisterAdmin(ctx context.Context, username, email, password string) (db.Admin, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return db.Admin{}, err
	}

	user, err := s.queries.CreateAdmin(ctx, db.CreateAdminParams{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		RoleID:       1,
		IsActive:     true,
	})
	if err != nil {
		return db.Admin{}, err
	}

	return user, nil
}

// SetEmailVerification sets the verification code and expiry for a user.
func (a *Service) SetEmailVerification(ctx context.Context, userID int32, code string, expiry time.Time) error {
	return a.queries.SetAdminEmailVerification(ctx, db.SetAdminEmailVerificationParams{
		ID:                    userID,
		VerificationCode:      sql.NullString{Valid: code != "", String: code},
		VerificationExpiresAt: sql.NullTime{Valid: true, Time: expiry},
	})
}

// VerifyEmailCode checks the code and marks the email as verified if valid and not expired.
func (a *Service) VerifyEmailCode(ctx context.Context, email, code string) (bool, error) {
	admin, err := a.queries.GetAdminByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	if admin.EmailVerified {
		return false, nil // Already verified
	}
	if admin.VerificationCode.String != code {
		return false, nil // Invalid code
	}
	if !admin.VerificationExpiresAt.Valid || admin.VerificationExpiresAt.Time.Before(time.Now()) {
		return false, nil // Expired
	}
	// Mark as verified and clear code
	err = a.queries.MarkAdminEmailVerified(ctx, db.MarkAdminEmailVerifiedParams{
		ID:            admin.ID,
		EmailVerified: true,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) Login(ctx context.Context, emailOrUsername, password string) (string, error) {
	// Try to find user by email first, then by username
	userByEmail, err := s.queries.GetUserByEmail(ctx, sql.NullString{String: emailOrUsername, Valid: true})
	if err != nil {
		// If not found by email, try by username
		userByUsername, err := s.queries.GetUserByUsername(ctx, emailOrUsername)
		if err != nil {
			return "", ErrInvalidCredentials
		}
		// Convert to common struct for consistent handling
		if !userByUsername.IsActive.Bool {
			return "", ErrUserInactive
		}

		if err := bcrypt.CompareHashAndPassword([]byte(userByUsername.PasswordHash), []byte(password)); err != nil {
			return "", ErrInvalidCredentials
		}

		permissions, err := s.queries.GetUserPermissions(ctx, userByUsername.ID)
		if err != nil {
			log.Println("Error getting user permissions:", err)
			return "", err
		}

		token, err := jwt.GenerateToken(
			int(userByUsername.ID),
			userByUsername.Email.String,
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
	if !userByEmail.IsActive.Bool {
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
		userByEmail.Email.String,
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

func(s *Service) Logout(ctx context.Context, token string, expiry time.Duration) error {
	claims, err := jwt.ParseToken(token, s.jwtSecret)
	if err != nil {
		return err
	}

	remainingTime := time.Until(claims.ExpiresAt.Time)
	if remainingTime > 0 {
		// add to blacklist until token expires
		err := s.redis.Set(ctx, fmt.Sprintf("jwt:blacklist:%s", token), "1", remainingTime)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
    exists, err := s.redis.Exists(ctx, fmt.Sprintf("jwt:blacklist:%s", token))
    return exists, err
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
	s.queries.UpdateUser(ctx, params)
	
	// Invalidate cache
	cacheKey := fmt.Sprintf("user:%d", params.ID)
	s.redis.Delete(ctx, cacheKey)
	
	return db.User{}, nil
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

func (s *Service) GetUserByID(ctx context.Context, id int32) (db.GetUserByIDRow, error) {
	cacheKey := fmt.Sprintf("user:%d", id)
	
	// Try cache first
	if cachedUser, err := s.redis.Get(ctx, cacheKey); err == nil {
		var user db.GetUserByIDRow
		err = json.Unmarshal([]byte(cachedUser), &user)
		if err != nil {
			return db.GetUserByIDRow{}, err
		}
		return user, nil
	}
	
	// Fetch from database if no cache
	user, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		return db.GetUserByIDRow{}, err
	}
	
	// Cache for 30 minutes
    jsonUser, _ := json.Marshal(user)
    s.redis.Set(ctx, cacheKey, jsonUser, 30*time.Minute)
	return user, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (db.GetUserByEmailRow, error) {
	cacheKey := fmt.Sprintf("user:email:%s", email)
	
	// Try cache first
	if cachedUser, err := s.redis.Get(ctx, cacheKey); err == nil {
		var user db.GetUserByEmailRow
		err = json.Unmarshal([]byte(cachedUser), &user)
		if err != nil {
			return db.GetUserByEmailRow{}, err
		}
		return user, nil
	}
	
	// Fetch from database if no cache
	user, err := s.queries.GetUserByEmail(ctx, sql.NullString{String: email, Valid: true})
	if err != nil {
		return db.GetUserByEmailRow{}, err
	}
	
	// Cache for 30 minutes
    jsonUser, _ := json.Marshal(user)
    s.redis.Set(ctx, cacheKey, jsonUser, 30*time.Minute)
	return user, nil
}

func (s *Service) GetUserByUsername(ctx context.Context, username string) (db.GetUserByUsernameRow, error) {
	cachedKey := fmt.Sprintf("user_by_username:%s", username)
	if cachedUser, err := s.redis.Get(ctx, cachedKey); err == nil {
		var user db.GetUserByUsernameRow
		err = json.Unmarshal([]byte(cachedUser), &user)
		if err != nil {
			return db.GetUserByUsernameRow{}, err
		}
		return user, nil
	}
	
	// Fetch from database if no cache
	user, err := s.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return db.GetUserByUsernameRow{}, err
	}
	
	// Cache for 30 minutes
    jsonUser, _ := json.Marshal(user)
    s.redis.Set(ctx, cachedKey, jsonUser, 30*time.Minute)
	return user, nil
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
func (s *Service) LogUserActivity(ctx context.Context, userID int, entityID int32, action, details, entityType, ip, userAgent string) error {
	_, err := s.queries.LogUserActivity(ctx, db.LogUserActivityParams{
		UserID:     int32(userID),
		Action:     action,
		Details:    json.RawMessage(details),
		EntityID:   entityID,
		EntityType: entityType,
		IpAddress:  sql.NullString{Valid: true, String: ip},
		UserAgent:  sql.NullString{Valid: true, String: userAgent},
	})
	return err
}

func (s *Service) LogLogin(ctx context.Context, username, email string, ip, userAgent string, success bool, errorReason string) error {
	err := s.queries.LogLoginHistory(ctx, db.LogLoginHistoryParams{
		Username:    username,
		Email:       sql.NullString{String: email, Valid: email != ""},
		IpAddress:   sql.NullString{String: ip, Valid: ip != ""},
		UserAgent:   sql.NullString{String: userAgent, Valid: userAgent != ""},
		ErrorReason: sql.NullString{String: errorReason, Valid: errorReason != ""},
	})
	return err
}
