package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	FindByPhone(ctx context.Context, phone string) (*User, error)
	Update(ctx context.Context, user *User) error
	UpdateVerificationStatus(ctx context.Context, email string, isVerified bool) error
	UpdateLastLogin(ctx context.Context, userID primitive.ObjectID) error
}

// RefreshTokenRepository defines the interface for refresh token data access
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *RefreshToken) error
	FindByToken(ctx context.Context, token string) (*RefreshToken, error)
	Revoke(ctx context.Context, token string) error
	RevokeAllForUser(ctx context.Context, userID primitive.ObjectID) error
}

// OTPRepository defines the interface for OTP storage (Redis)
type OTPRepository interface {
	Store(ctx context.Context, email, code string, ttl int) error
	Get(ctx context.Context, email string) (string, error)
	Delete(ctx context.Context, email string) error
	Exists(ctx context.Context, email string) (bool, error)
}

// RateLimiterRepository defines the interface for rate limiting
type RateLimiterRepository interface {
	CheckRateLimit(ctx context.Context, key string, limit int, window int) (bool, error)
	IncrementCounter(ctx context.Context, key string, window int) error
}
