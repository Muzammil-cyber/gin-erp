package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	domain "github.com/muzammil-cyber/gin-erp/internal/domain/auth"
	"github.com/muzammil-cyber/gin-erp/pkg/utils"
)

type AuthHandler struct {
	authUseCase domain.AuthUseCase
}

func NewAuthHandler(authUseCase domain.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with email, phone, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.RegisterRequest true "Registration request"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorMessageResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authUseCase.Register(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrUserAlreadyExists, domain.ErrPhoneAlreadyExists:
			utils.ErrorResponse(c, http.StatusConflict, err)
		case domain.ErrInvalidPhoneFormat, domain.ErrInvalidRole:
			utils.ErrorResponse(c, http.StatusBadRequest, err)
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, err)
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, gin.H{
		"message": "Registration successful. Please verify your email with the OTP sent.",
		"user":    user.ToUserInfo(),
	})
}

// Login handles user login
// @Summary Login user
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Login request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorMessageResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	authResponse, err := h.authUseCase.Login(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			utils.ErrorResponse(c, http.StatusUnauthorized, err)
		case domain.ErrUserNotVerified:
			utils.ErrorResponse(c, http.StatusForbidden, err)
		case domain.ErrUserInactive:
			utils.ErrorResponse(c, http.StatusForbidden, err)
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, err)
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, authResponse)
}

// VerifyOTP handles OTP verification
// @Summary Verify OTP
// @Description Verify OTP code sent to user's email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.VerifyOTPRequest true "OTP verification request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /auth/verify-otp [post]
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req domain.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorMessageResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.authUseCase.VerifyOTP(c.Request.Context(), &req); err != nil {
		switch err {
		case domain.ErrInvalidOTP, domain.ErrOTPNotFound:
			utils.ErrorResponse(c, http.StatusBadRequest, err)
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, err)
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"message": "Email verified successfully. You can now login.",
	})
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req domain.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorMessageResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	authResponse, err := h.authUseCase.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrInvalidToken, domain.ErrExpiredToken, domain.ErrRevokedToken:
			utils.ErrorResponse(c, http.StatusUnauthorized, err)
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, err)
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, authResponse)
}

// ResendOTP handles OTP resend
// @Summary Resend OTP
// @Description Resend OTP code to user's email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Email"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /auth/resend-otp [post]
func (h *AuthHandler) ResendOTP(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorMessageResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.authUseCase.ResendOTP(c.Request.Context(), req.Email); err != nil {
		switch err {
		case domain.ErrUserNotFound:
			utils.ErrorResponse(c, http.StatusNotFound, err)
		case domain.ErrOTPAlreadySent:
			utils.ErrorResponse(c, http.StatusTooManyRequests, err)
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, err)
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"message": "OTP sent successfully. Please check your email.",
	})
}

// GetProfile gets the current user's profile
// @Summary Get user profile
// @Description Get authenticated user's profile
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, domain.ErrUnauthorized)
		return
	}

	user, err := h.authUseCase.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, user.ToUserInfo())
}
