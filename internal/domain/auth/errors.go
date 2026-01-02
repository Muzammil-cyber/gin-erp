package domain

import "errors"

var (
	// Authentication errors
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrPhoneAlreadyExists = errors.New("user with this phone number already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserNotVerified    = errors.New("user email not verified")
	ErrUserInactive       = errors.New("user account is inactive")

	// Token errors
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
	ErrRevokedToken = errors.New("token has been revoked")

	// OTP errors
	ErrInvalidOTP     = errors.New("invalid OTP code")
	ErrExpiredOTP     = errors.New("OTP has expired")
	ErrOTPNotFound    = errors.New("OTP not found")
	ErrOTPAlreadySent = errors.New("OTP already sent, please wait before requesting again")

	// Validation errors
	ErrInvalidPhoneFormat = errors.New("invalid Pakistani phone number format (expected: +923xxxxxxxxx)")
	ErrInvalidRole        = errors.New("invalid role")
	ErrWeakPassword       = errors.New("password must be at least 8 characters")

	// Rate limiting errors
	ErrRateLimitExceeded = errors.New("rate limit exceeded, please try again later")

	// Generic errors
	ErrInternalServer = errors.New("internal server error")
	ErrBadRequest     = errors.New("bad request")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden: insufficient permissions")
)
