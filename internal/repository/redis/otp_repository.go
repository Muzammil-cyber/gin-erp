package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/muzammil-cyber/gin-erp/pkg/database/redis"
)

type OTPRepositoryImpl struct {
	redisClient *redis.Client
}

func NewOTPRepository(redisClient *redis.Client) *OTPRepositoryImpl {
	return &OTPRepositoryImpl{
		redisClient: redisClient,
	}
}

// Store stores an OTP in Redis with TTL
func (r *OTPRepositoryImpl) Store(ctx context.Context, email, code string, ttlMinutes int) error {
	key := fmt.Sprintf("otp:%s", email)
	return r.redisClient.Set(ctx, key, code, time.Duration(ttlMinutes)*time.Minute)
}

// Get retrieves an OTP from Redis
func (r *OTPRepositoryImpl) Get(ctx context.Context, email string) (string, error) {
	key := fmt.Sprintf("otp:%s", email)
	return r.redisClient.Get(ctx, key)
}

// Delete removes an OTP from Redis
func (r *OTPRepositoryImpl) Delete(ctx context.Context, email string) error {
	key := fmt.Sprintf("otp:%s", email)
	return r.redisClient.Delete(ctx, key)
}

// Exists checks if an OTP exists in Redis
func (r *OTPRepositoryImpl) Exists(ctx context.Context, email string) (bool, error) {
	key := fmt.Sprintf("otp:%s", email)
	return r.redisClient.Exists(ctx, key)
}
