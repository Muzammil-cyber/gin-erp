package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muzammil-cyber/gin-erp/internal/delivery/http/handler"
	"github.com/muzammil-cyber/gin-erp/internal/delivery/http/middleware"
	"github.com/muzammil-cyber/gin-erp/internal/delivery/http/routes"
	domain "github.com/muzammil-cyber/gin-erp/internal/domain/auth"
	mongoRepo "github.com/muzammil-cyber/gin-erp/internal/repository/mongodb"
	redisRepo "github.com/muzammil-cyber/gin-erp/internal/repository/redis"
	"github.com/muzammil-cyber/gin-erp/internal/service"
	"github.com/muzammil-cyber/gin-erp/internal/usecase"
	"github.com/muzammil-cyber/gin-erp/pkg/database/mongodb"
	"github.com/muzammil-cyber/gin-erp/pkg/database/redis"
	"github.com/muzammil-cyber/gin-erp/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoginIntegration is an integration test for the login endpoint
// This test requires MongoDB and Redis to be running
// Run with: go test -v ./tests/integration -tags=integration
func TestLoginIntegration(t *testing.T) {
	// Skip if not in integration test mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test environment
	ctx := context.Background()

	// Initialize test database connections
	mongoClient, err := mongodb.NewMongoDBClient(
		"mongodb://localhost:27017",
		"test_pakistani_erp",
		10,
		5,
	)
	require.NoError(t, err, "Failed to connect to MongoDB")
	defer mongoClient.Close(ctx)

	redisClient, err := redis.NewRedisClient(
		"localhost",
		"6379",
		"",
		1, // Use DB 1 for tests
	)
	require.NoError(t, err, "Failed to connect to Redis")
	defer redisClient.Close()

	// Initialize repositories
	userRepo := mongoRepo.NewUserRepository(mongoClient)
	refreshTokenRepo := mongoRepo.NewRefreshTokenRepository(mongoClient)
	otpRepo := redisRepo.NewOTPRepository(redisClient)
	rateLimiterRepo := redisRepo.NewRateLimiterRepository(redisClient)

	// Create indexes
	err = userRepo.CreateIndexes(ctx)
	require.NoError(t, err)

	// Initialize services
	emailService := service.NewEmailService("", 587, "", "", "")
	jwtUtil := utils.NewJWTUtil("test-secret-key", 15*time.Minute, 168*time.Hour)

	// Initialize use case
	authUseCase := usecase.NewAuthUseCase(
		userRepo,
		refreshTokenRepo,
		otpRepo,
		emailService,
		jwtUtil,
		6,
		5,
	)

	// Initialize handler
	authHandler := handler.NewAuthHandler(authUseCase)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.TraceIDMiddleware())

	routerConfig := &routes.RouterConfig{
		AuthHandler:    authHandler,
		JWTUtil:        jwtUtil,
		RateLimiter:    rateLimiterRepo,
		LoginRateLimit: 5,
	}
	routes.SetupRoutes(router, routerConfig)

	// Cleanup: Delete test user if exists
	testEmail := "integration.test@example.com"
	existingUser, _ := userRepo.FindByEmail(ctx, testEmail)
	if existingUser != nil {
		// In production code, you would have a Delete method
		// For now, we'll just update the user
	}

	// Test 1: Register a new user
	t.Run("Register User", func(t *testing.T) {
		registerReq := domain.RegisterRequest{
			Email:     testEmail,
			Phone:     "+923001234567",
			Password:  "TestPassword123",
			FirstName: "Integration",
			LastName:  "Test",
			Role:      domain.RoleCustomer,
		}

		body, _ := json.Marshal(registerReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response utils.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.NotEmpty(t, response.TraceID)
	})

	// Test 2: Verify OTP (mock - we'll manually set verification)
	t.Run("Verify User", func(t *testing.T) {
		// Manually verify the user for testing
		err := userRepo.UpdateVerificationStatus(ctx, testEmail, true)
		require.NoError(t, err)
	})

	// Test 3: Login with correct credentials
	t.Run("Login Success", func(t *testing.T) {
		loginReq := domain.LoginRequest{
			Email:    testEmail,
			Password: "TestPassword123",
		}

		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response utils.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.NotNil(t, response.Data)

		// Verify response structure
		data := response.Data.(map[string]interface{})
		assert.NotEmpty(t, data["access_token"])
		assert.NotEmpty(t, data["refresh_token"])
		assert.NotNil(t, data["user"])

		// Store tokens for next test
		accessToken := data["access_token"].(string)
		assert.NotEmpty(t, accessToken)
	})

	// Test 4: Login with incorrect credentials
	t.Run("Login Failure - Wrong Password", func(t *testing.T) {
		loginReq := domain.LoginRequest{
			Email:    testEmail,
			Password: "WrongPassword",
		}

		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response utils.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.False(t, response.Success)
		assert.NotEmpty(t, response.Error)
	})

	// Test 5: Login with non-existent user
	t.Run("Login Failure - User Not Found", func(t *testing.T) {
		loginReq := domain.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "SomePassword",
		}

		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response utils.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.False(t, response.Success)
	})

	// Test 6: Get user profile with valid token
	t.Run("Get Profile Success", func(t *testing.T) {
		// First login to get a valid token
		loginReq := domain.LoginRequest{
			Email:    testEmail,
			Password: "TestPassword123",
		}

		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var loginResponse utils.Response
		json.Unmarshal(w.Body.Bytes(), &loginResponse)
		data := loginResponse.Data.(map[string]interface{})
		accessToken := data["access_token"].(string)

		// Now get profile
		req = httptest.NewRequest(http.MethodGet, "/api/v1/auth/profile", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response utils.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.NotNil(t, response.Data)
	})

	// Test 7: Get profile without token
	t.Run("Get Profile Failure - No Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/profile", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response utils.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.False(t, response.Success)
	})
}
