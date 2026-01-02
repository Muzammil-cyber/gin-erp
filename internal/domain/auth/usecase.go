package domain

import "context"

// AuthUseCase defines the business logic interface for authentication
type AuthUseCase interface {
	Register(ctx context.Context, req *RegisterRequest) (*User, error)
	Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error)
	VerifyOTP(ctx context.Context, req *VerifyOTPRequest) error
	RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*AuthResponse, error)
	ResendOTP(ctx context.Context, email string) error
	GetUserByID(ctx context.Context, userID string) (*User, error)
}

// EmailService defines the interface for sending emails
type EmailService interface {
	SendOTP(ctx context.Context, email, code string) error
}
