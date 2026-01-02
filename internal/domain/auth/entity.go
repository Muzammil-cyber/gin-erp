package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Role represents user roles in the system
type Role string

const (
	RoleAdmin          Role = "admin"
	RoleCustomer       Role = "customer"
	RoleFinanceManager Role = "finance_manager"
	RoleManager        Role = "manager"
)

// User represents a user entity in the system
type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email       string             `json:"email" bson:"email"`
	Phone       string             `json:"phone" bson:"phone"`
	Password    string             `json:"-" bson:"password"`
	FirstName   string             `json:"first_name" bson:"first_name"`
	LastName    string             `json:"last_name" bson:"last_name"`
	Role        Role               `json:"role" bson:"role"`
	IsVerified  bool               `json:"is_verified" bson:"is_verified"`
	IsActive    bool               `json:"is_active" bson:"is_active"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	LastLoginAt *time.Time         `json:"last_login_at,omitempty" bson:"last_login_at,omitempty"`
}

// RefreshToken represents a refresh token entity
type RefreshToken struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	Token     string             `json:"token" bson:"token"`
	ExpiresAt time.Time          `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	IsRevoked bool               `json:"is_revoked" bson:"is_revoked"`
}

// OTP represents an OTP entity for verification
type OTP struct {
	Email     string    `json:"email"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Role      Role   `json:"role" binding:"required,oneof=admin customer finance_manager manager"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// VerifyOTPRequest represents the OTP verification request
type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

// RefreshTokenRequest represents the refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	User         *UserInfo `json:"user"`
}

// UserInfo represents user information in response
type UserInfo struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Role       Role      `json:"role"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
}

// ToUserInfo converts User to UserInfo
func (u *User) ToUserInfo() *UserInfo {
	return &UserInfo{
		ID:         u.ID.Hex(),
		Email:      u.Email,
		Phone:      u.Phone,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		Role:       u.Role,
		IsVerified: u.IsVerified,
		CreatedAt:  u.CreatedAt,
	}
}

// IsValidRole checks if the role is valid
func IsValidRole(role Role) bool {
	switch role {
	case RoleAdmin, RoleCustomer, RoleFinanceManager, RoleManager:
		return true
	}
	return false
}
