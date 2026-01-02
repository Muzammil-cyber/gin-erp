package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	domain "github.com/muzammil-cyber/gin-erp/internal/domain/auth"
	"github.com/muzammil-cyber/gin-erp/internal/usecase"
	"github.com/muzammil-cyber/gin-erp/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock repositories
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	args := m.Called(ctx, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateVerificationStatus(ctx context.Context, email string, isVerified bool) error {
	args := m.Called(ctx, email, isVerified)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, userID primitive.ObjectID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) Revoke(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID primitive.ObjectID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

type MockOTPRepository struct {
	mock.Mock
}

func (m *MockOTPRepository) Store(ctx context.Context, email, code string, ttl int) error {
	args := m.Called(ctx, email, code, ttl)
	return args.Error(0)
}

func (m *MockOTPRepository) Get(ctx context.Context, email string) (string, error) {
	args := m.Called(ctx, email)
	return args.String(0), args.Error(1)
}

func (m *MockOTPRepository) Delete(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockOTPRepository) Exists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendOTP(ctx context.Context, email, code string) error {
	args := m.Called(ctx, email, code)
	return args.Error(0)
}

// Table-driven tests for Register
func TestAuthUseCase_Register(t *testing.T) {
	tests := []struct {
		name          string
		request       *domain.RegisterRequest
		mockSetup     func(*MockUserRepository, *MockOTPRepository, *MockEmailService)
		expectedError error
	}{
		{
			name: "Successful registration",
			request: &domain.RegisterRequest{
				Email:     "test@example.com",
				Phone:     "+923001234567",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
				Role:      domain.RoleCustomer,
			},
			mockSetup: func(userRepo *MockUserRepository, otpRepo *MockOTPRepository, emailService *MockEmailService) {
				userRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(nil, domain.ErrUserNotFound)
				userRepo.On("FindByPhone", mock.Anything, "+923001234567").Return(nil, domain.ErrUserNotFound)
				userRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
				otpRepo.On("Store", mock.Anything, "test@example.com", mock.AnythingOfType("string"), 5).Return(nil)
				emailService.On("SendOTP", mock.Anything, "test@example.com", mock.AnythingOfType("string")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Invalid phone number format",
			request: &domain.RegisterRequest{
				Email:     "test@example.com",
				Phone:     "invalid-phone",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
				Role:      domain.RoleCustomer,
			},
			mockSetup: func(userRepo *MockUserRepository, otpRepo *MockOTPRepository, emailService *MockEmailService) {
				// No mocks needed as validation fails before DB calls
			},
			expectedError: domain.ErrInvalidPhoneFormat,
		},
		{
			name: "User already exists",
			request: &domain.RegisterRequest{
				Email:     "existing@example.com",
				Phone:     "+923001234567",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
				Role:      domain.RoleCustomer,
			},
			mockSetup: func(userRepo *MockUserRepository, otpRepo *MockOTPRepository, emailService *MockEmailService) {
				existingUser := &domain.User{
					ID:    primitive.NewObjectID(),
					Email: "existing@example.com",
				}
				userRepo.On("FindByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil)
			},
			expectedError: domain.ErrUserAlreadyExists,
		},
		{
			name: "Invalid role",
			request: &domain.RegisterRequest{
				Email:     "test@example.com",
				Phone:     "+923001234567",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
				Role:      "invalid_role",
			},
			mockSetup: func(userRepo *MockUserRepository, otpRepo *MockOTPRepository, emailService *MockEmailService) {
				// No mocks needed
			},
			expectedError: domain.ErrInvalidRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			userRepo := new(MockUserRepository)
			refreshTokenRepo := new(MockRefreshTokenRepository)
			otpRepo := new(MockOTPRepository)
			emailService := new(MockEmailService)
			jwtUtil := utils.NewJWTUtil("test-secret", 15*time.Minute, 168*time.Hour)

			tt.mockSetup(userRepo, otpRepo, emailService)

			// Create use case
			uc := usecase.NewAuthUseCase(userRepo, refreshTokenRepo, otpRepo, emailService, jwtUtil, 6, 5)

			// Execute
			user, err := uc.Register(context.Background(), tt.request)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.request.Email, user.Email)
				assert.Equal(t, "+923001234567", user.Phone)
				assert.False(t, user.IsVerified)
				assert.True(t, user.IsActive)
			}

			userRepo.AssertExpectations(t)
			otpRepo.AssertExpectations(t)
			emailService.AssertExpectations(t)
		})
	}
}

// Table-driven tests for Login
func TestAuthUseCase_Login(t *testing.T) {
	hashedPassword, _ := utils.HashPassword("password123")
	userID := primitive.NewObjectID()

	tests := []struct {
		name          string
		request       *domain.LoginRequest
		mockSetup     func(*MockUserRepository, *MockRefreshTokenRepository)
		expectedError error
	}{
		{
			name: "Successful login",
			request: &domain.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(userRepo *MockUserRepository, tokenRepo *MockRefreshTokenRepository) {
				user := &domain.User{
					ID:         userID,
					Email:      "test@example.com",
					Password:   hashedPassword,
					Role:       domain.RoleCustomer,
					IsVerified: true,
					IsActive:   true,
				}
				userRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(user, nil)
				tokenRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.RefreshToken")).Return(nil)
				userRepo.On("UpdateLastLogin", mock.Anything, userID).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Invalid credentials - wrong password",
			request: &domain.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(userRepo *MockUserRepository, tokenRepo *MockRefreshTokenRepository) {
				user := &domain.User{
					Email:    "test@example.com",
					Password: hashedPassword,
				}
				userRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(user, nil)
			},
			expectedError: domain.ErrInvalidCredentials,
		},
		{
			name: "User not found",
			request: &domain.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			mockSetup: func(userRepo *MockUserRepository, tokenRepo *MockRefreshTokenRepository) {
				userRepo.On("FindByEmail", mock.Anything, "nonexistent@example.com").Return(nil, domain.ErrUserNotFound)
			},
			expectedError: domain.ErrInvalidCredentials,
		},
		{
			name: "User not verified",
			request: &domain.LoginRequest{
				Email:    "unverified@example.com",
				Password: "password123",
			},
			mockSetup: func(userRepo *MockUserRepository, tokenRepo *MockRefreshTokenRepository) {
				user := &domain.User{
					Email:      "unverified@example.com",
					Password:   hashedPassword,
					IsVerified: false,
					IsActive:   true,
				}
				userRepo.On("FindByEmail", mock.Anything, "unverified@example.com").Return(user, nil)
			},
			expectedError: domain.ErrUserNotVerified,
		},
		{
			name: "User inactive",
			request: &domain.LoginRequest{
				Email:    "inactive@example.com",
				Password: "password123",
			},
			mockSetup: func(userRepo *MockUserRepository, tokenRepo *MockRefreshTokenRepository) {
				user := &domain.User{
					Email:      "inactive@example.com",
					Password:   hashedPassword,
					IsVerified: true,
					IsActive:   false,
				}
				userRepo.On("FindByEmail", mock.Anything, "inactive@example.com").Return(user, nil)
			},
			expectedError: domain.ErrUserInactive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			userRepo := new(MockUserRepository)
			refreshTokenRepo := new(MockRefreshTokenRepository)
			otpRepo := new(MockOTPRepository)
			emailService := new(MockEmailService)
			jwtUtil := utils.NewJWTUtil("test-secret", 15*time.Minute, 168*time.Hour)

			tt.mockSetup(userRepo, refreshTokenRepo)

			// Create use case
			uc := usecase.NewAuthUseCase(userRepo, refreshTokenRepo, otpRepo, emailService, jwtUtil, 6, 5)

			// Execute
			authResponse, err := uc.Login(context.Background(), tt.request)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, authResponse)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, authResponse)
				assert.NotEmpty(t, authResponse.AccessToken)
				assert.NotEmpty(t, authResponse.RefreshToken)
				assert.NotNil(t, authResponse.User)
			}

			userRepo.AssertExpectations(t)
			refreshTokenRepo.AssertExpectations(t)
		})
	}
}

// Test VerifyOTP
func TestAuthUseCase_VerifyOTP(t *testing.T) {
	tests := []struct {
		name          string
		request       *domain.VerifyOTPRequest
		mockSetup     func(*MockOTPRepository, *MockUserRepository)
		expectedError error
	}{
		{
			name: "Valid OTP",
			request: &domain.VerifyOTPRequest{
				Email: "test@example.com",
				Code:  "123456",
			},
			mockSetup: func(otpRepo *MockOTPRepository, userRepo *MockUserRepository) {
				otpRepo.On("Get", mock.Anything, "test@example.com").Return("123456", nil)
				userRepo.On("UpdateVerificationStatus", mock.Anything, "test@example.com", true).Return(nil)
				otpRepo.On("Delete", mock.Anything, "test@example.com").Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Invalid OTP",
			request: &domain.VerifyOTPRequest{
				Email: "test@example.com",
				Code:  "wrong-code",
			},
			mockSetup: func(otpRepo *MockOTPRepository, userRepo *MockUserRepository) {
				otpRepo.On("Get", mock.Anything, "test@example.com").Return("123456", nil)
			},
			expectedError: domain.ErrInvalidOTP,
		},
		{
			name: "OTP not found",
			request: &domain.VerifyOTPRequest{
				Email: "test@example.com",
				Code:  "123456",
			},
			mockSetup: func(otpRepo *MockOTPRepository, userRepo *MockUserRepository) {
				otpRepo.On("Get", mock.Anything, "test@example.com").Return("", errors.New("not found"))
			},
			expectedError: domain.ErrOTPNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(MockUserRepository)
			refreshTokenRepo := new(MockRefreshTokenRepository)
			otpRepo := new(MockOTPRepository)
			emailService := new(MockEmailService)
			jwtUtil := utils.NewJWTUtil("test-secret", 15*time.Minute, 168*time.Hour)

			tt.mockSetup(otpRepo, userRepo)

			uc := usecase.NewAuthUseCase(userRepo, refreshTokenRepo, otpRepo, emailService, jwtUtil, 6, 5)
			err := uc.VerifyOTP(context.Background(), tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			otpRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
		})
	}
}
