package auth

import (
	"context"
	"database/sql"
	"errors"
	db "herp/db/sqlc"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type mockQuerier struct {
	db.Queries
	createAdminParams db.CreateAdminParams
	createAdminResp   db.Admin
	createAdminErr    error

	getAdminByEmailResp db.GetAdminByEmailRow
	getAdminByEmailErr  error

	setAdminEmailVerificationCalled bool
	markAdminEmailVerifiedCalled    bool

	// logLoginAttemptCalled bool
}

func (m *mockQuerier) CreateAdmin(ctx context.Context, params db.CreateAdminParams) (db.Admin, error) {
	m.createAdminParams = params
	return m.createAdminResp, m.createAdminErr
}
func (m *mockQuerier) SetAdminEmailVerification(ctx context.Context, params db.SetAdminEmailVerificationParams) error {
	m.setAdminEmailVerificationCalled = true
	return nil
}
func (m *mockQuerier) GetAdminByEmail(ctx context.Context, email string) (db.GetAdminByEmailRow, error) {
	return m.getAdminByEmailResp, m.getAdminByEmailErr
}
func (m *mockQuerier) MarkAdminEmailVerified(ctx context.Context, params db.MarkAdminEmailVerifiedParams) error {
	m.markAdminEmailVerifiedCalled = true
	return nil
}

func TestRegisterAdmin(t *testing.T) {
	mockQ := &mockQuerier{
		createAdminResp: db.Admin{
			ID:           1,
			Username:     "admin",
			Email:        "admin@example.com",
			FirstName:    "Admin",
			LastName:     "User",
			PasswordHash: "",
			RoleID:       1,
			IsActive:     true,
		},
	}
	svc := &Service{queries: mockQ}

	admin, err := svc.RegisterAdmin(context.Background(), "admin", "admin@example.com", "password", "Admin", "User")
	require.NoError(t, err)
	assert.Equal(t, "admin", admin.Username)
	assert.Equal(t, "admin@example.com", admin.Email)
	assert.True(t, mockQ.createAdminParams.IsActive)
	assert.Equal(t, int32(1), mockQ.createAdminParams.RoleID)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(mockQ.createAdminParams.PasswordHash), []byte("password")))
}

func TestSetEmailVerification(t *testing.T) {
	mockQ := &mockQuerier{}
	svc := &Service{queries: mockQ}
	expiry := time.Now().Add(10 * time.Minute)
	err := svc.SetEmailVerification(context.Background(), 42, "verification-token", expiry)
	require.NoError(t, err)
	assert.True(t, mockQ.setAdminEmailVerificationCalled)
}

func TestVerifyEmailCode_Success(t *testing.T) {
	expiry := time.Now().Add(10 * time.Minute)
	mockQ := &mockQuerier{
		getAdminByEmailResp: db.GetAdminByEmailRow{
			ID:                    1,
			EmailVerified:         false,
			VerificationCode:      sql.NullString{String: "code123", Valid: true},
			VerificationExpiresAt: sql.NullTime{Time: expiry, Valid: true},
		},
	}
	svc := &Service{queries: mockQ}
	ok, err := svc.VerifyEmailCode(context.Background(), "admin@example.com", "code123")
	require.NoError(t, err)
	assert.True(t, ok)
	assert.True(t, mockQ.markAdminEmailVerifiedCalled)
}

func TestVerifyEmailCode_AlreadyVerified(t *testing.T) {
	mockQ := &mockQuerier{
		getAdminByEmailResp: db.GetAdminByEmailRow{
			ID:            1,
			EmailVerified: true,
		},
	}
	svc := &Service{queries: mockQ}
	ok, err := svc.VerifyEmailCode(context.Background(), "admin@example.com", "any")
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestVerifyEmailCode_InvalidCode(t *testing.T) {
	expiry := time.Now().Add(10 * time.Minute)
	mockQ := &mockQuerier{
		getAdminByEmailResp: db.GetAdminByEmailRow{
			ID:                    1,
			EmailVerified:         false,
			VerificationCode:      sql.NullString{String: "rightcode", Valid: true},
			VerificationExpiresAt: sql.NullTime{Time: expiry, Valid: true},
		},
	}
	svc := &Service{queries: mockQ}
	ok, err := svc.VerifyEmailCode(context.Background(), "admin@example.com", "wrongcode")
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestVerifyEmailCode_Expired(t *testing.T) {
	expiry := time.Now().Add(-10 * time.Minute)
	mockQ := &mockQuerier{
		getAdminByEmailResp: db.GetAdminByEmailRow{
			ID:                    1,
			EmailVerified:         false,
			VerificationCode:      sql.NullString{String: "code123", Valid: true},
			VerificationExpiresAt: sql.NullTime{Time: expiry, Valid: true},
		},
	}
	svc := &Service{queries: mockQ}
	ok, err := svc.VerifyEmailCode(context.Background(), "admin@example.com", "code123")
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestVerifyEmailCode_DBError(t *testing.T) {
	mockQ := &mockQuerier{
		getAdminByEmailErr: errors.New("db error"),
	}
	svc := &Service{queries: mockQ}
	ok, err := svc.VerifyEmailCode(context.Background(), "admin@example.com", "code123")
	assert.Error(t, err)
	assert.False(t, ok)
}

