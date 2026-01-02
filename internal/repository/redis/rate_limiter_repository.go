package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/muzammil-cyber/gin-erp/pkg/database/redis"
)

type RateLimiterRepositoryImpl struct {
	redisClient *redis.Client
}

func NewRateLimiterRepository(redisClient *redis.Client) *RateLimiterRepositoryImpl {
	return &RateLimiterRepositoryImpl{
		redisClient: redisClient,
	}
}

// CheckRateLimit checks if the rate limit has been exceeded
func (r *RateLimiterRepositoryImpl) CheckRateLimit(ctx context.Context, key string, limit int, windowSeconds int) (bool, error) {
	rateLimitKey := fmt.Sprintf("rate_limit:%s", key)

	// Get current count
	count, err := r.redisClient.GetClient().Get(ctx, rateLimitKey).Int()
	if err != nil && err.Error() != "redis: nil" {
		return false, err
	}

	// Check if limit exceeded
	if count >= limit {
		return false, nil // Rate limit exceeded
	}

	return true, nil // Within rate limit
}

// IncrementCounter increments the rate limit counter
func (r *RateLimiterRepositoryImpl) IncrementCounter(ctx context.Context, key string, windowSeconds int) error {
	rateLimitKey := fmt.Sprintf("rate_limit:%s", key)

	// Increment counter
	count, err := r.redisClient.Increment(ctx, rateLimitKey)
	if err != nil {
		return err
	}

	// Set expiration on first increment
	if count == 1 {
		if err := r.redisClient.Expire(ctx, rateLimitKey, time.Duration(windowSeconds)*time.Second); err != nil {
			return err
		}
	}

	return nil
}
