package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muzammil-cyber/gin-erp/internal/container"
	"github.com/muzammil-cyber/gin-erp/internal/delivery/http/middleware"
	"github.com/muzammil-cyber/gin-erp/internal/delivery/http/routes"

	_ "github.com/muzammil-cyber/gin-erp/docs" // swagger docs
)

// @title Pakistani ERP System API
// @version 1.0
// @description A production-ready ERP system API built with Gin framework following DDD principles
// @description
// @description Features:
// @description - JWT-based authentication with token rotation
// @description - Role-based access control (Admin, Customer, Finance Manager, Manager)
// @description - Pakistani phone number validation (+923xxxxxxxxx)
// @description - Email OTP verification with 5-minute TTL
// @description - Redis-based rate limiting
// @description - MongoDB for data storage
// @description - Redis for caching and OTP storage

// @contact.name Muzammil Loya
// @contact.url https://github.com/muzammilloya
// @contact.email muzammil@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Initialize dependency container
	c, err := container.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer c.Close()

	// Set Gin mode
	if c.Config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.TraceIDMiddleware())
	router.Use(middleware.CORSMiddleware(
		c.Config.CORS.AllowedOrigins,
		c.Config.CORS.AllowedMethods,
		c.Config.CORS.AllowedHeaders,
	))
	router.Use(middleware.RateLimiterMiddleware(
		c.RateLimiterRepo,
		c.Config.RateLimit.RequestsPerMinute,
		60, // 60 seconds window
	))

	// Setup routes
	routerConfig := &routes.RouterConfig{
		AuthHandler:    c.AuthHandler,
		JWTUtil:        c.JWTUtil,
		RateLimiter:    c.RateLimiterRepo,
		LoginRateLimit: c.Config.RateLimit.LoginRequestsPerMinute,
	}
	routes.SetupRoutes(router, routerConfig)

	// Create HTTP server
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", c.Config.App.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Server starting on port %s", c.Config.App.Port)
		log.Printf("📍 Environment: %s", c.Config.App.Env)
		log.Printf("🔗 Server URL: http://localhost:%s", c.Config.App.Port)
		log.Printf("🏥 Health check: http://localhost:%s/health", c.Config.App.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Graceful shutdown with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("✓ Server exited gracefully")
}
