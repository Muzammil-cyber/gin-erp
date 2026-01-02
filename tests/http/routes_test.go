package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/muzammil-cyber/gin-erp/internal/delivery/http/handler"
	"github.com/muzammil-cyber/gin-erp/internal/delivery/http/routes"
	domain "github.com/muzammil-cyber/gin-erp/internal/domain/auth"
	"github.com/muzammil-cyber/gin-erp/pkg/utils"
)

// MockAuthUseCase is a mock implementation of AuthUseCase
type MockAuthUseCase struct {
	mock.Mock
}

func (m *MockAuthUseCase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockAuthUseCase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AuthResponse), args.Error(1)
}

func (m *MockAuthUseCase) VerifyOTP(ctx context.Context, req *domain.VerifyOTPRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockAuthUseCase) RefreshToken(ctx context.Context, req *domain.RefreshTokenRequest) (*domain.AuthResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AuthResponse), args.Error(1)
}

func (m *MockAuthUseCase) ResendOTP(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockAuthUseCase) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// MockRateLimiter is a mock implementation of RateLimiterRepository
type MockRateLimiter struct {
	mock.Mock
}

func (m *MockRateLimiter) CheckRateLimit(ctx context.Context, key string, limit int, window int) (bool, error) {
	args := m.Called(ctx, key, limit, window)
	return args.Bool(0), args.Error(1)
}

func (m *MockRateLimiter) IncrementCounter(ctx context.Context, key string, window int) error {
	args := m.Called(ctx, key, window)
	return args.Error(0)
}

// setupTestRouter creates a test router with mocked dependencies
func setupTestRouter() (*gin.Engine, *MockAuthUseCase, *MockRateLimiter) {
	gin.SetMode(gin.TestMode)

	mockAuthUseCase := new(MockAuthUseCase)
	mockRateLimiter := new(MockRateLimiter)

	authHandler := handler.NewAuthHandler(mockAuthUseCase)
	jwtUtil := utils.NewJWTUtil("test-secret-key", 15*60, 7*24*60*60)

	router := gin.New()
	config := &routes.RouterConfig{
		AuthHandler:    authHandler,
		JWTUtil:        jwtUtil,
		RateLimiter:    mockRateLimiter,
		LoginRateLimit: 5,
	}

	routes.SetupRoutes(router, config)

	return router, mockAuthUseCase, mockRateLimiter
}

func TestHealthCheckRoute(t *testing.T) {
	router, _, _ := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "ok", data["status"])
	assert.Equal(t, "Pakistani ERP System", data["service"])
}

func TestSwaggerRoute(t *testing.T) {
	router, _, _ := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/swagger/index.html", nil)
	router.ServeHTTP(w, req)

	// Swagger should return 200 or 301 (redirect) or 404 (if docs not generated)
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusMovedPermanently || w.Code == http.StatusNotFound)
}

func TestRegisterRoute(t *testing.T) {
	router, mockUseCase, _ := setupTestRouter()

	t.Run("Successful Registration", func(t *testing.T) {
		mockUseCase.On("Register", mock.Anything, mock.Anything).Return(&domain.User{
			Email:      "test@example.com",
			Phone:      "+923001234567",
			FirstName:  "Test",
			LastName:   "User",
			Role:       domain.RoleCustomer,
			IsVerified: false,
		}, nil).Once()

		requestBody := `{
			"email": "test@example.com",
			"phone": "+923001234567",
			"password": "Password123",
			"first_name": "Test",
			"last_name": "User",
			"role": "customer"
		}`

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		requestBody := `{"invalid": "data"}`

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLoginRoute(t *testing.T) {
	router, mockUseCase, mockRateLimiter := setupTestRouter()

	t.Run("Successful Login", func(t *testing.T) {
		mockRateLimiter.On("CheckRateLimit", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(true, nil).Once()
		mockRateLimiter.On("IncrementCounter", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockUseCase.On("Login", mock.Anything, mock.Anything).Return(&domain.AuthResponse{
			AccessToken:  "access_token",
			RefreshToken: "refresh_token",
		}, nil).Once()

		requestBody := `{
			"email": "test@example.com",
			"password": "Password123"
		}`

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))
	})

	t.Run("Rate Limit Exceeded", func(t *testing.T) {
		mockRateLimiter.On("CheckRateLimit", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(false, nil).Once()

		requestBody := `{
			"email": "test@example.com",
			"password": "Password123"
		}`

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
	})
}

func TestProtectedRoutes(t *testing.T) {
	router, _, _ := setupTestRouter()

	tests := []struct {
		name           string
		endpoint       string
		method         string
		expectedStatus int
	}{
		{
			name:           "Profile Without Auth",
			endpoint:       "/api/v1/auth/profile",
			method:         "GET",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Admin Endpoint Without Auth",
			endpoint:       "/api/v1/admin/users",
			method:         "GET",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Finance Endpoint Without Auth",
			endpoint:       "/api/v1/finance/reports",
			method:         "GET",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Manager Endpoint Without Auth",
			endpoint:       "/api/v1/manager/dashboard",
			method:         "GET",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.endpoint, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
