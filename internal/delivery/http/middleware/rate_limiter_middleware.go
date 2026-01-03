package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	domain "github.com/muzammil-cyber/gin-erp/internal/domain/auth"
	"github.com/muzammil-cyber/gin-erp/pkg/utils"
)

// RateLimiterMiddleware implements Redis-based rate limiting
func RateLimiterMiddleware(rateLimiter domain.RateLimiterRepository, limit int, windowSeconds int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client IP or user ID for rate limiting
		key := c.ClientIP()

		// If user is authenticated, use user ID instead
		if userID, exists := c.Get("user_id"); exists {
			key = fmt.Sprintf("user:%s", userID.(string))
		} else {
			key = fmt.Sprintf("ip:%s", key)
		}

		// Check rate limit
		allowed, err := rateLimiter.CheckRateLimit(c.Request.Context(), key, limit, windowSeconds)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, domain.ErrInternalServer)
			c.Abort()
			return
		}

		if !allowed {
			utils.ErrorResponse(c, http.StatusTooManyRequests, domain.ErrRateLimitExceeded)
			c.Abort()
			return
		}

		// Increment counter
		if err := rateLimiter.IncrementCounter(c.Request.Context(), key, windowSeconds); err != nil {
			// Log error but don't fail the request
			log.Printf("Rate limiter increment error: %v", err)
		}

		c.Next()
	}
}

// LoginRateLimiterMiddleware is a specific rate limiter for login endpoints
func LoginRateLimiterMiddleware(rateLimiter domain.RateLimiterRepository, limit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use IP + endpoint for login rate limiting
		key := fmt.Sprintf("login:ip:%s", c.ClientIP())

		// Check rate limit (60 seconds window)
		allowed, err := rateLimiter.CheckRateLimit(c.Request.Context(), key, limit, 60)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, domain.ErrInternalServer)
			c.Abort()
			return
		}

		if !allowed {
			utils.ErrorMessageResponse(c, http.StatusTooManyRequests, "Too many login attempts. Please try again later.")
			c.Abort()
			return
		}

		// Increment counter
		if err := rateLimiter.IncrementCounter(c.Request.Context(), key, 60); err != nil {
			// Log error but don't fail the request
			log.Printf("Rate limiter increment error: %v", err)
		}

		c.Next()
	}
}
