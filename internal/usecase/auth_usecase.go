package usecase

import (
	"context"
	"dormitory_management/internal/config"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/entity/helper"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/internal/domain/response"
	"dormitory_management/internal/domain/usecase"
	"errors"
	"fmt"
	"time"

	"dormitory_management/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// Claims is the custom JWT claims
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// authUseCase implements the usecase.AuthUseCase interface
type authUseCase struct {
	repos  repository.Repositories
	logger logger.Logger
	config *config.Config
}

// NewAuthUseCase creates a new auth use case
func NewAuthUseCase(repos repository.Repositories, logger logger.Logger) usecase.AuthUseCase {
	cfg := config.LoadConfig()
	return &authUseCase{
		repos:  repos,
		logger: logger,
		config: cfg,
	}
}

// Register registers a new user
func (uc *authUseCase) Register(ctx context.Context, userData *entity.UserRegister) response.StatusResponse {
	// Check if email already exists
	_, err := uc.repos.User().GetByEmail(ctx, userData.Email)
	if err == nil {
		return response.BadRequest("Email already exists")
	}

	// Create user with student role
	user := &entity.User{
		FullName: userData.FullName,
		Email:    userData.Email,
		Password: userData.Password,
		Role:     entity.UserRoleStudent,
		Gender:   entity.UserGenderOther,
		Status:   entity.UserStatusActive,
	}

	if err := uc.repos.Auth().Register(ctx, user); err != nil {
		uc.logger.Error("Failed to register user", zap.Error(err))
		return response.InternalServerError("Failed to register user")
	}

	// Generate tokens
	accessToken, refreshToken, err := uc.GenerateTokens(ctx, user.ID)
	if err != nil {
		uc.logger.Error("Failed to generate tokens", zap.Error(err))
		return response.InternalServerError("Failed to generate tokens")
	}

	return response.Success(map[string]interface{}{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, 1)
}

// Login authenticates a user and returns tokens
func (uc *authUseCase) Login(ctx context.Context, loginData *entity.UserLogin) response.StatusResponse {
	user, err := uc.repos.Auth().Login(ctx, loginData.Email, loginData.Password)
	if err != nil {
		uc.logger.Error("Login failed", zap.Error(err))
		return response.Unauthorized("Invalid email or password")
	}

	// Generate tokens
	accessToken, refreshToken, err := uc.GenerateTokens(ctx, user.ID)
	if err != nil {
		uc.logger.Error("Failed to generate tokens", zap.Error(err))
		return response.InternalServerError("Failed to generate tokens")
	}

	return response.Success(map[string]interface{}{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, 1)
}

// RefreshToken refreshes an access token
func (uc *authUseCase) RefreshToken(ctx context.Context, refreshToken string) response.StatusResponse {
	// Verify refresh token in Redis
	userID, err := uc.repos.Auth().GetToken(ctx, refreshToken)
	if err != nil {
		uc.logger.Error("Invalid refresh token", zap.Error(err))
		return response.Unauthorized("Invalid refresh token")
	}

	// Get user
	user, err := uc.repos.User().GetByID(ctx, userID)
	if err != nil {
		uc.logger.Error("User not found", zap.Error(err))
		return response.Unauthorized("User not found")
	}

	// Generate new tokens
	accessToken, newRefreshToken, err := uc.GenerateTokens(ctx, user.ID)
	if err != nil {
		uc.logger.Error("Failed to generate tokens", zap.Error(err))
		return response.InternalServerError("Failed to generate tokens")
	}

	// Invalidate old refresh token
	if err := uc.repos.Auth().InvalidateToken(ctx, user.ID, refreshToken); err != nil {
		uc.logger.Error("Failed to invalidate old refresh token", zap.Error(err))
		// Continue anyway
	}

	return response.Success(map[string]interface{}{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	}, 1)
}

// Logout invalidates tokens
func (uc *authUseCase) Logout(ctx context.Context, userID uint, token string) response.StatusResponse {
	if err := uc.repos.Auth().InvalidateToken(ctx, userID, token); err != nil {
		uc.logger.Error("Failed to invalidate token", zap.Error(err))
		return response.InternalServerError("Failed to logout")
	}

	return response.Success("Logged out successfully", 0)
}

// GenerateTokens generates access and refresh tokens
func (uc *authUseCase) GenerateTokens(ctx context.Context, userID uint) (string, string, error) {
	// Get user
	user, err := uc.repos.User().GetByID(ctx, userID)
	if err != nil {
		return "", "", errors.New("user not found")
	}

	// Generate access token
	accessTokenClaims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(uc.config.JWT.AccessExpiresIn) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "dormitory-management",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(uc.config.JWT.Secret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %w", err)
	}

	// Generate refresh token (just a random string in practice)
	refreshTokenString := helper.GenerateRandomString(32)

	// Store refresh token in Redis
	if err := uc.repos.Auth().StoreToken(ctx, user.ID, refreshTokenString, uc.config.JWT.RefreshExpiresIn); err != nil {
		return "", "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return accessTokenString, refreshTokenString, nil
}
