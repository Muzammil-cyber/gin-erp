package usecase

import (
	"context"
	"log"
	"time"

	domain "github.com/muzammil-cyber/gin-erp/internal/domain/auth"
	"github.com/muzammil-cyber/gin-erp/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthUseCaseImpl struct {
	userRepo         domain.UserRepository
	refreshTokenRepo domain.RefreshTokenRepository
	otpRepo          domain.OTPRepository
	emailService     domain.EmailService
	jwtUtil          *utils.JWTUtil
	otpLength        int
	otpExpireMinutes int
}

func NewAuthUseCase(
	userRepo domain.UserRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	otpRepo domain.OTPRepository,
	emailService domain.EmailService,
	jwtUtil *utils.JWTUtil,
	otpLength int,
	otpExpireMinutes int,
) domain.AuthUseCase {
	return &AuthUseCaseImpl{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		otpRepo:          otpRepo,
		emailService:     emailService,
		jwtUtil:          jwtUtil,
		otpLength:        otpLength,
		otpExpireMinutes: otpExpireMinutes,
	}
}

// Register registers a new user
func (uc *AuthUseCaseImpl) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error) {
	// Validate Pakistani phone number
	normalizedPhone := utils.NormalizePakistaniPhone(req.Phone)
	if !utils.ValidatePakistaniPhone(normalizedPhone) {
		return nil, domain.ErrInvalidPhoneFormat
	}

	// Validate role
	if !domain.IsValidRole(req.Role) {
		return nil, domain.ErrInvalidRole
	}

	// Check if user already exists
	existingUser, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// Check if phone already exists
	existingPhone, err := uc.userRepo.FindByPhone(ctx, normalizedPhone)
	if err == nil && existingPhone != nil {
		return nil, domain.ErrPhoneAlreadyExists
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &domain.User{
		Email:      req.Email,
		Phone:      normalizedPhone,
		Password:   hashedPassword,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Role:       req.Role,
		IsVerified: false,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate and send OTP
	otp, err := utils.GenerateOTP(uc.otpLength)
	if err != nil {
		return nil, err
	}

	if err := uc.otpRepo.Store(ctx, req.Email, otp, uc.otpExpireMinutes); err != nil {
		return nil, err
	}

	if err := uc.emailService.SendOTP(ctx, req.Email, otp); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns tokens
func (uc *AuthUseCaseImpl) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	// Find user by email
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Check password
	if !utils.CheckPassword(user.Password, req.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	// Check if user is verified
	if !user.IsVerified {
		return nil, domain.ErrUserNotVerified
	}

	// Check if user is active
	if !user.IsActive {
		return nil, domain.ErrUserInactive
	}

	// Generate access token
	accessToken, err := uc.jwtUtil.GenerateAccessToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := uc.jwtUtil.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Store refresh token in database
	refreshTokenEntity := &domain.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(uc.jwtUtil.GetRefreshTokenExpiration()),
		IsRevoked: false,
	}

	if err := uc.refreshTokenRepo.Create(ctx, refreshTokenEntity); err != nil {
		return nil, err
	}

	// Update last login
	if err := uc.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Log error but don't fail the login
		log.Printf("Failed to update last login for user %s: %v", user.ID.Hex(), err)
	}

	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToUserInfo(),
	}, nil
}

// VerifyOTP verifies the OTP code
func (uc *AuthUseCaseImpl) VerifyOTP(ctx context.Context, req *domain.VerifyOTPRequest) error {
	// Get OTP from Redis
	storedOTP, err := uc.otpRepo.Get(ctx, req.Email)
	if err != nil {
		return domain.ErrOTPNotFound
	}

	// Verify OTP
	if storedOTP != req.Code {
		return domain.ErrInvalidOTP
	}

	// Update user verification status
	if err := uc.userRepo.UpdateVerificationStatus(ctx, req.Email, true); err != nil {
		return err
	}

	// Delete OTP from Redis
	if err := uc.otpRepo.Delete(ctx, req.Email); err != nil {
		// Log error but don't fail
		log.Printf("Failed to delete OTP for email %s: %v", req.Email, err)
	}

	return nil
}

// RefreshToken refreshes the access token using a refresh token
func (uc *AuthUseCaseImpl) RefreshToken(ctx context.Context, req *domain.RefreshTokenRequest) (*domain.AuthResponse, error) {
	// Validate refresh token
	claims, err := uc.jwtUtil.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	// Check if token exists and is not revoked
	refreshToken, err := uc.refreshTokenRepo.FindByToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	if refreshToken.IsRevoked {
		return nil, domain.ErrRevokedToken
	}

	// Check if token is expired
	if time.Now().After(refreshToken.ExpiresAt) {
		return nil, domain.ErrExpiredToken
	}

	// Get user
	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// Generate new access token
	accessToken, err := uc.jwtUtil.GenerateAccessToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	// Generate new refresh token
	newRefreshToken, err := uc.jwtUtil.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Revoke old refresh token
	if err := uc.refreshTokenRepo.Revoke(ctx, req.RefreshToken); err != nil {
		return nil, err
	}

	// Store new refresh token
	newRefreshTokenEntity := &domain.RefreshToken{
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(uc.jwtUtil.GetRefreshTokenExpiration()),
		IsRevoked: false,
	}

	if err := uc.refreshTokenRepo.Create(ctx, newRefreshTokenEntity); err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         user.ToUserInfo(),
	}, nil
}

// ResendOTP resends the OTP to the user
func (uc *AuthUseCaseImpl) ResendOTP(ctx context.Context, email string) error {
	// Check if user exists
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return domain.ErrUserNotFound
	}

	// Check if user is already verified
	if user.IsVerified {
		return domain.ErrBadRequest
	}

	// Check if OTP already exists (prevent spam)
	exists, err := uc.otpRepo.Exists(ctx, email)
	if err != nil {
		return err
	}

	if exists {
		return domain.ErrOTPAlreadySent
	}

	// Generate and send new OTP
	otp, err := utils.GenerateOTP(uc.otpLength)
	if err != nil {
		return err
	}

	if err := uc.otpRepo.Store(ctx, email, otp, uc.otpExpireMinutes); err != nil {
		return err
	}

	return uc.emailService.SendOTP(ctx, email, otp)
}

// GetUserByID retrieves a user by ID
func (uc *AuthUseCaseImpl) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	return uc.userRepo.FindByID(ctx, id)
}
