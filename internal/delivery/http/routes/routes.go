package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/muzammil-cyber/gin-erp/internal/delivery/http/handler"
	"github.com/muzammil-cyber/gin-erp/internal/delivery/http/middleware"
	domain "github.com/muzammil-cyber/gin-erp/internal/domain/auth"
	"github.com/muzammil-cyber/gin-erp/pkg/utils"
)

type RouterConfig struct {
	AuthHandler    *handler.AuthHandler
	JWTUtil        *utils.JWTUtil
	RateLimiter    domain.RateLimiterRepository
	LoginRateLimit int
}

// SetupRoutes sets up all routes
func SetupRoutes(router *gin.Engine, config *RouterConfig) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		utils.SuccessResponse(c, 200, gin.H{
			"status":  "ok",
			"service": "Pakistani ERP System",
		})
	})

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Public auth routes
		auth := v1.Group("/auth")
		{
			// Apply login rate limiter to login endpoint
			auth.POST("/register", config.AuthHandler.Register)
			auth.POST("/login",
				middleware.LoginRateLimiterMiddleware(config.RateLimiter, config.LoginRateLimit),
				config.AuthHandler.Login,
			)
			auth.POST("/verify-otp", config.AuthHandler.VerifyOTP)
			auth.POST("/refresh-token", config.AuthHandler.RefreshToken)
			auth.POST("/resend-otp", config.AuthHandler.ResendOTP)
		}

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(config.JWTUtil))
		{
			protected.GET("/auth/profile", config.AuthHandler.GetProfile)

			// Admin only routes
			admin := protected.Group("")
			admin.Use(middleware.RoleMiddleware(domain.RoleAdmin))
			{
				// Add admin routes here
				admin.GET("/admin/users", func(c *gin.Context) {
					utils.SuccessResponse(c, 200, gin.H{
						"message": "Admin users endpoint",
					})
				})
			}

			// Finance manager routes
			finance := protected.Group("")
			finance.Use(middleware.RoleMiddleware(domain.RoleAdmin, domain.RoleFinanceManager))
			{
				// Add finance routes here
				finance.GET("/finance/reports", func(c *gin.Context) {
					utils.SuccessResponse(c, 200, gin.H{
						"message": "Finance reports endpoint",
					})
				})
			}

			// Manager routes
			manager := protected.Group("")
			manager.Use(middleware.RoleMiddleware(domain.RoleAdmin, domain.RoleManager))
			{
				// Add manager routes here
				manager.GET("/manager/dashboard", func(c *gin.Context) {
					utils.SuccessResponse(c, 200, gin.H{
						"message": "Manager dashboard endpoint",
					})
				})
			}
		}
	}
}
