package container

import (
	"context"
	"fmt"
	"log"

	"github.com/muzammil-cyber/gin-erp/internal/config"
	"github.com/muzammil-cyber/gin-erp/internal/delivery/http/handler"
	mongoRepo "github.com/muzammil-cyber/gin-erp/internal/repository/mongodb"
	redisRepo "github.com/muzammil-cyber/gin-erp/internal/repository/redis"
	"github.com/muzammil-cyber/gin-erp/internal/service"
	"github.com/muzammil-cyber/gin-erp/internal/usecase"
	"github.com/muzammil-cyber/gin-erp/pkg/database/mongodb"
	"github.com/muzammil-cyber/gin-erp/pkg/database/redis"
	"github.com/muzammil-cyber/gin-erp/pkg/utils"
)

// Container holds all dependencies
type Container struct {
	Config  *config.Config
	MongoDB *mongodb.Client
	Redis   *redis.Client
	JWTUtil *utils.JWTUtil

	// Repositories
	UserRepo         *mongoRepo.UserRepositoryImpl
	RefreshTokenRepo *mongoRepo.RefreshTokenRepositoryImpl
	OTPRepo          *redisRepo.OTPRepositoryImpl
	RateLimiterRepo  *redisRepo.RateLimiterRepositoryImpl

	// Services
	EmailService *service.EmailServiceImpl

	// Use cases
	AuthUseCase *usecase.AuthUseCaseImpl

	// Handlers
	AuthHandler *handler.AuthHandler
}

// NewContainer initializes all dependencies
func NewContainer() (*Container, error) {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	log.Printf("Initializing %s in %s mode...", cfg.App.Name, cfg.App.Env)

	// Initialize MongoDB
	log.Println("Connecting to MongoDB...")
	mongoClient, err := mongodb.NewMongoDBClient(
		cfg.MongoDB.URI,
		cfg.MongoDB.Database,
		cfg.MongoDB.MaxPoolSize,
		cfg.MongoDB.MinPoolSize,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	log.Println("✓ MongoDB connected successfully")

	// Initialize Redis
	log.Println("Connecting to Redis...")
	redisClient, err := redis.NewRedisClient(
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	log.Println("✓ Redis connected successfully")

	// Initialize JWT utility
	jwtUtil := utils.NewJWTUtil(
		cfg.JWT.Secret,
		cfg.GetAccessTokenDuration(),
		cfg.GetRefreshTokenDuration(),
	)

	// Initialize repositories
	userRepo := mongoRepo.NewUserRepository(mongoClient)
	refreshTokenRepo := mongoRepo.NewRefreshTokenRepository(mongoClient)
	otpRepo := redisRepo.NewOTPRepository(redisClient)
	rateLimiterRepo := redisRepo.NewRateLimiterRepository(redisClient)

	// Create indexes
	log.Println("Creating database indexes...")
	ctx := context.Background()
	if err := userRepo.CreateIndexes(ctx); err != nil {
		log.Printf("Warning: Failed to create user indexes: %v", err)
	}
	if err := refreshTokenRepo.CreateIndexes(ctx); err != nil {
		log.Printf("Warning: Failed to create refresh token indexes: %v", err)
	}
	log.Println("✓ Database indexes created")

	// Initialize services
	emailService := service.NewEmailService(
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.Username,
		cfg.SMTP.Password,
		cfg.SMTP.From,
	)

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(
		userRepo,
		refreshTokenRepo,
		otpRepo,
		emailService,
		jwtUtil,
		cfg.OTP.Length,
		cfg.OTP.ExpireMinutes,
	)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUseCase)

	log.Println("✓ All dependencies initialized successfully")

	return &Container{
		Config:           cfg,
		MongoDB:          mongoClient,
		Redis:            redisClient,
		JWTUtil:          jwtUtil,
		UserRepo:         userRepo,
		RefreshTokenRepo: refreshTokenRepo,
		OTPRepo:          otpRepo,
		RateLimiterRepo:  rateLimiterRepo,
		EmailService:     emailService,
		AuthUseCase:      authUseCase.(*usecase.AuthUseCaseImpl),
		AuthHandler:      authHandler,
	}, nil
}

// Close closes all connections
func (c *Container) Close() error {
	log.Println("Closing connections...")

	ctx := context.Background()

	if err := c.MongoDB.Close(ctx); err != nil {
		log.Printf("Error closing MongoDB: %v", err)
	}

	if err := c.Redis.Close(); err != nil {
		log.Printf("Error closing Redis: %v", err)
	}

	log.Println("✓ All connections closed")
	return nil
}
