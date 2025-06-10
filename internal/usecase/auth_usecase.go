package usecase

import (
	"context"
	"dormitory_management/internal/common"
	"dormitory_management/internal/config"
	"dormitory_management/internal/delivery/mq/tasks"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/internal/domain/response"
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/internal/helper"
	"dormitory_management/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Claims is the custom JWT claims
type Claims struct {
	UserID       uint   `json:"user_id"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	TokenVersion string `json:"token_version"`
	jwt.RegisteredClaims
}

// authUseCase implements the usecase.AuthUseCase interface
type authUseCase struct {
	repos       repository.Repositories
	logger      logger.Logger
	config      *config.Config
	asynqClient *asynq.Client
}

// NewAuthUseCase creates a new auth use case
func NewAuthUseCase(repos repository.Repositories, logger logger.Logger, asynqClient *asynq.Client) usecase.AuthUseCase {
	cfg := config.LoadConfig()
	return &authUseCase{
		repos:       repos,
		logger:      logger,
		config:      cfg,
		asynqClient: asynqClient,
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
		Phone:    userData.Phone,
		Role:     entity.UserRoleStudent,
		Gender:   entity.UserGenderOther,
		Status:   entity.UserStatusActive,
	}
	otpCode := &entity.OtpCode{
		IsUsed:         false,
		OtpCode:        common.GenerateCode(6),
		UserID:         &user.ID,
		IdentifierType: entity.IdentifierTypeEmail.String(),
		Identifier:     user.Email,
		OtpType:        entity.OtpTypeVerifyEmail.String(),
		ExpiresAt:      common.GetExpireTime(15),
	}

	if err := uc.repos.Auth().Register(ctx, user, otpCode); err != nil {
		uc.logger.Error("Failed to register user", zap.Error(err))
		return response.InternalServerError("Failed to register user")
	}

	jsonPayload, err := json.Marshal(&entity.SendMailVerifyAccount{
		UserID: user.ID,
		Email:  user.Email,
		Token:  uuid.New().String(),
	})
	if err != nil {
		uc.logger.Error("Failed to marshal payload", zap.Error(err))
		return response.InternalServerError("Failed send mail verify account")
	}
	task := asynq.NewTask(string(tasks.TypeSendEmailVerifyAccount), jsonPayload)
	uc.asynqClient.EnqueueContext(ctx, task)

	return response.Success(user, 1)
}

// Me returns the current user
func (uc *authUseCase) Me(ctx context.Context, userID uint) response.StatusResponse {
	// Get user from cache
	user, err := uc.repos.Auth().GetUserCache(ctx, userID)
	if errors.Is(err, redis.Nil) {
		uc.logger.Info("User not found in cache, fetching from database")
		user, err := uc.repos.User().GetByID(ctx, userID)
		if err != nil {
			uc.logger.Error("User not found", zap.Error(err))
			return response.Unauthorized("User not found")
		}

		err = uc.repos.Auth().SetUserCache(ctx, user, uc.config.JWT.AccessExpiresIn)
		if err != nil {
			uc.logger.Error("Failed to set user cache", zap.Error(err))
			return response.InternalServerError("Failed to set user cache")
		}

		return response.Success(user, 1)
	}

	if err != nil {
		uc.logger.Error("User not found", zap.Error(err))
		return response.Unauthorized("User not found")
	}

	// user, err := uc.repos.User().GetByID(ctx, userID)
	// if err != nil {
	// 	uc.logger.Error("User not found", zap.Error(err))
	// 	return response.Unauthorized("User not found")
	// }

	return response.Success(user, 1)
}

// Login authenticates a user and returns tokens
func (uc *authUseCase) Login(ctx context.Context, loginData *entity.UserLogin) response.StatusResponse {
	user, err := uc.repos.Auth().Login(ctx, loginData.Email, loginData.Password)
	if err != nil {
		uc.logger.Error("Login failed", zap.Error(err))
		return response.Unauthorized("Invalid email or password")
	}

	if !user.IsVerify {
		uc.logger.Error("Unverified user")
		return response.Forbidden("Unverified user")
	}

	// Generate tokens
	accessToken, refreshToken, err := uc.GenerateTokens(ctx, user.ID)
	if err != nil {
		uc.logger.Error("Failed to generate tokens", zap.Error(err))
		return response.InternalServerError("Failed to generate tokens")
	}

	secretKey := uc.config.Server.SecretKey

	encryptAccessToken, err := helper.Encrypt(accessToken, secretKey)
	if err != nil {
		uc.logger.Error("Failed encrypt access token", zap.Error(err))
		return response.InternalServerError("Failed encrypt access token")
	}

	encryptRefreshToken, err := helper.Encrypt(refreshToken, secretKey)
	if err != nil {
		uc.logger.Error("Failed encrypt refresh token", zap.Error(err))
		return response.InternalServerError("Failed encrypt refresh token")
	}

	return response.Success(map[string]interface{}{
		"user":          user,
		"access_token":  encryptAccessToken,
		"refresh_token": encryptRefreshToken,
	}, 1)
}

// RefreshToken refreshes an access token
func (uc *authUseCase) RefreshToken(ctx context.Context, refreshToken string) response.StatusResponse {
	decryptRefreshToken, err := helper.Decrypt(refreshToken, uc.config.Server.SecretKey)
	if err != nil {
		uc.logger.Error("Failed to decrypt token", zap.Error(err))
		return response.Unauthorized(err.Error())
	}

	token, err := jwt.ParseWithClaims(decryptRefreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(uc.config.JWT.RefreshSecret), nil
	})
	if err != nil {
		uc.logger.Error("Failed to parse refresh token", zap.Error(err))
		return response.Unauthorized("Invalid refresh token")
	}
	if !token.Valid {
		uc.logger.Error("Invalid refresh token", zap.Error(err))
		return response.Unauthorized("Invalid refresh token")
	}
	// Extract the claims from the token
	claims, ok := token.Claims.(*Claims)
	if !ok {
		uc.logger.Error("Failed to extract claims from token", zap.Error(err))
		return response.Unauthorized("Invalid refresh token")
	}

	// Verify refresh token in Redis
	tokenVersion, err := uc.repos.Auth().CheckTokenVersion(ctx, entity.RefreshToken, claims.UserID)
	if err != nil {
		uc.logger.Error("Failed to check token version", zap.Error(err))
		return response.Unauthorized("Invalid refresh token")
	}
	if tokenVersion != claims.TokenVersion {
		uc.logger.Error("Failed to check token version", zap.Error(err))
		return response.Unauthorized("Invalid refresh token")
	}

	// Get user
	user, err := uc.repos.User().GetByID(ctx, claims.UserID)
	if err != nil {
		uc.logger.Error("User not found", zap.Error(err))
		return response.Unauthorized("User not found")
	}

	// Generate new tokens
	accessToken, err := uc.createAccessToken(ctx, user)
	if err != nil {
		uc.logger.Error("Failed to create access token", zap.Error(err))
		return response.InternalServerError("Failed to create access token")
	}

	encryptAccessToken, err := helper.Encrypt(accessToken, uc.config.Server.SecretKey)
	if err != nil {
		uc.logger.Error("Failed encrypt access token", zap.Error(err))
		return response.InternalServerError("Failed encrypt access token")
	}

	return response.Success(map[string]interface{}{
		"access_token": encryptAccessToken,
	}, 1)
}

// Logout invalidates tokens
func (uc *authUseCase) Logout(ctx context.Context, userID uint) response.StatusResponse {
	if err := uc.repos.Auth().InvalidateToken(ctx, entity.AccessToken, userID); err != nil {
		uc.logger.Error("Failed to invalidate token", zap.Error(err))
		return response.InternalServerError("Failed to logout")
	}

	if err := uc.repos.Auth().InvalidateToken(ctx, entity.RefreshToken, userID); err != nil {
		uc.logger.Error("Failed to invalidate token", zap.Error(err))
		return response.InternalServerError("Failed to logout")
	}

	if err := uc.repos.Auth().DeleteUserCache(ctx, userID); err != nil {
		uc.logger.Error("Failed to delete user cache", zap.Error(err))
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

	accessToken, err := uc.createAccessToken(ctx, user)
	if err != nil {
		return "", "", fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, err := uc.createRefreshToken(ctx, user)
	if err != nil {
		return "", "", fmt.Errorf("failed to create refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (uc *authUseCase) createToken(tokenClaims Claims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(uc.config.JWT.AccessSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return accessTokenString, nil
}

// createAccessToken creates an access token
func (uc *authUseCase) createAccessToken(ctx context.Context, user *entity.User) (string, error) {
	expiresIn := uc.config.JWT.AccessExpiresIn
	uuid := uuid.New()
	accessTokenClaims := Claims{
		UserID:       user.ID,
		Email:        user.Email,
		Role:         string(user.Role),
		TokenVersion: uuid.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresIn) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "dormitory-management",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	accessToken, err := uc.createToken(accessTokenClaims)
	if err != nil {
		return "", err
	}

	if err := uc.repos.Auth().SetCacheTokenVersion(ctx, entity.AccessToken, user.ID, string(uuid.String()), expiresIn); err != nil {
		uc.logger.Error("Failed to set token version", zap.Error(err))
		return "", fmt.Errorf("failed to set token version: %w", err)
	}

	if err := uc.repos.Auth().SetUserCache(ctx, user, expiresIn); err != nil {
		uc.logger.Error("Failed to set user cache", zap.Error(err))
		return "", fmt.Errorf("failed to set user cache: %w", err)
	}

	return accessToken, nil
}

// createRefreshToken creates a refresh token
func (uc *authUseCase) createRefreshToken(ctx context.Context, user *entity.User) (string, error) {
	uuid := uuid.New()
	refreshTokenClaims := Claims{
		UserID:       user.ID,
		Email:        user.Email,
		Role:         string(user.Role),
		TokenVersion: uuid.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(uc.config.JWT.RefreshExpiresIn) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "dormitory-management",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	refreshToken, err := uc.createToken(refreshTokenClaims)
	if err != nil {
		return "", err
	}

	if err := uc.repos.Auth().SetCacheTokenVersion(ctx, entity.RefreshToken, user.ID, string(uuid.String()), uc.config.JWT.RefreshExpiresIn); err != nil {
		uc.logger.Error("Failed to set token version", zap.Error(err))
		return "", fmt.Errorf("failed to set token version: %w", err)
	}

	return refreshToken, nil
}
