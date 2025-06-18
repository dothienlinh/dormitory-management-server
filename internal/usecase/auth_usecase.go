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
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Claims struct {
	UserID       uint64 `json:"user_id"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	TokenVersion string `json:"token_version"`
	jwt.RegisteredClaims
}

type authUseCase struct {
	repos       repository.Repositories
	logger      logger.Logger
	config      *config.Config
	asynqClient *asynq.Client
}

func NewAuthUseCase(repos repository.Repositories, logger logger.Logger, asynqClient *asynq.Client) usecase.AuthUseCase {
	cfg := config.LoadConfig()
	return &authUseCase{
		repos:       repos,
		logger:      logger,
		config:      cfg,
		asynqClient: asynqClient,
	}
}

func (uc *authUseCase) Register(ctx context.Context, userData *entity.UserRegister) response.StatusResponse {
	if err := uc.repos.User().GetByEmail(ctx, &entity.User{Email: userData.Email}); err == nil {
		return response.BadRequest("Email already exists")
	}

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
		OtpType:        entity.OtpTypeVerifyAccount.String(),
		ExpiresAt:      common.GetExpireTime(15),
	}

	if err := uc.repos.Auth().Register(ctx, user, otpCode); err != nil {
		uc.logger.Error("Failed to register user", zap.Error(err))
		return response.InternalServerError("Failed to register user")
	}

	otpCodeEncrypt, err := helper.Encrypt(otpCode.OtpCode, uc.config.Server.SecretKey)
	if err != nil {
		uc.logger.Error("Failed encrypt otp code", zap.Error(err))
		return response.InternalServerError("Failed encrypt otp code")
	}

	jsonPayload, err := json.Marshal(&entity.SendMailVerifyAccount{
		UserID: user.ID,
		Email:  user.Email,
		Token:  otpCodeEncrypt,
	})
	if err != nil {
		uc.logger.Error("Failed to marshal payload", zap.Error(err))
		return response.InternalServerError("Failed send mail verify account")
	}
	task := asynq.NewTask(string(tasks.TypeSendEmailVerifyAccount), jsonPayload)
	uc.asynqClient.EnqueueContext(ctx, task)

	return response.Success(user, 1)
}

func (uc *authUseCase) Me(ctx context.Context, userID uint64) response.StatusResponse {
	user := &entity.User{
		Base: entity.Base{ID: userID},
	}
	if err := uc.repos.Auth().GetUserCache(ctx, user); errors.Is(err, redis.Nil) {
		uc.logger.Info("User not found in cache, fetching from database")
		if err := uc.repos.User().GetByID(ctx, user); err != nil {
			uc.logger.Error("User not found", zap.Error(err))
			return response.Unauthorized("User not found")
		}

		err = uc.repos.Auth().SetUserCache(ctx, user, uc.config.JWT.AccessExpiresIn)
		if err != nil {
			uc.logger.Error("Failed to set user cache", zap.Error(err))
			return response.InternalServerError("Failed to set user cache")
		}

		return response.Success(user, 1)
	} else if err != nil {
		uc.logger.Error("User not found", zap.Error(err))
		return response.Unauthorized("User not found")
	}

	return response.Success(user, 1)
}

func (uc *authUseCase) Login(ctx context.Context, loginData *entity.UserLogin) response.StatusResponse {
	user := &entity.User{
		Email:    loginData.Email,
		Password: loginData.Password,
	}
	if err := uc.repos.Auth().Login(ctx, user); err != nil {
		uc.logger.Error("Login failed", zap.Error(err))
		return response.Unauthorized(err.Error())
	}

	if !user.IsVerify {
		uc.logger.Error("Unverified user")
		return response.Forbidden("Unverified user")
	}

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

func (uc *authUseCase) RefreshToken(ctx context.Context, refreshToken string) response.StatusResponse {
	decryptRefreshToken, err := helper.Decrypt(refreshToken, uc.config.Server.SecretKey)
	if err != nil {
		uc.logger.Error("Failed to decrypt token", zap.Error(err))
		return response.Unauthorized(err.Error())
	}

	token, err := jwt.ParseWithClaims(decryptRefreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
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
	claims, ok := token.Claims.(*Claims)
	if !ok {
		uc.logger.Error("Failed to extract claims from token", zap.Error(err))
		return response.Unauthorized("Invalid refresh token")
	}

	tokenVersion, err := uc.repos.Auth().CheckTokenVersion(ctx, entity.RefreshToken, claims.UserID)
	if err != nil {
		uc.logger.Error("Failed to check token version", zap.Error(err))
		return response.Unauthorized("Invalid refresh token")
	}
	if tokenVersion != claims.TokenVersion {
		uc.logger.Error("Failed to check token version", zap.Error(err))
		return response.Unauthorized("Invalid refresh token")
	}

	user := &entity.User{
		Base: entity.Base{ID: claims.UserID},
	}
	if err := uc.repos.User().GetByID(ctx, user); err != nil {
		uc.logger.Error("User not found", zap.Error(err))
		return response.Unauthorized("User not found")
	}

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

func (uc *authUseCase) Logout(ctx context.Context, userID uint64) response.StatusResponse {
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

func (uc *authUseCase) GenerateTokens(ctx context.Context, userID uint64) (string, string, error) {
	user := &entity.User{
		Base: entity.Base{ID: userID},
	}
	if err := uc.repos.User().GetByID(ctx, user); err != nil {
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

func (uc *authUseCase) VerifyAccount(ctx context.Context, payload entity.VerifyAccount) response.StatusResponse {

	unescapeToken, err := url.QueryUnescape(payload.Token)
	if err != nil {
		uc.logger.Error("Failed to unescape token", zap.Error(err))
		return response.InternalServerError(err.Error())
	}

	tokenDecrypt, err := helper.Decrypt(unescapeToken, uc.config.Server.SecretKey)
	if err != nil {
		uc.logger.Error("Failed to decrypt token", zap.Error(err))
		return response.Unauthorized(err.Error())
	}

	otpCode := &entity.OtpCode{
		OtpCode:    tokenDecrypt,
		Identifier: payload.Email,
	}
	if err := uc.repos.OtpCode().FindCode(ctx, otpCode); err != nil {
		uc.logger.Error("Failed to find otp code", zap.Error(err))
		return response.Unauthorized(err.Error())
	}

	if otpCode.ID == 0 {
		return response.Unauthorized("Otp code not found")
	}

	if otpCode.IsUsed || otpCode.VerifiedAt != nil {
		return response.Unauthorized("Otp code used")
	}

	if otpCode.Identifier != payload.Email ||
		otpCode.IdentifierType != entity.IdentifierTypeEmail.String() ||
		otpCode.OtpType != entity.OtpTypeVerifyAccount.String() {
		return response.Unauthorized("Invalid otp code")
	}

	expiresAt, err := common.ParsedTime(otpCode.ExpiresAt)
	if err != nil {
		uc.logger.Error("Failed to parse expiresAt", zap.Error(err))
		return response.InternalServerError("Failed to verify OTP code")
	}
	if time.Now().After(expiresAt) {
		uc.logger.Error("OTP code expired")
		return response.BadRequest("OTP code expired")
	}

	user := &entity.User{
		Email: payload.Email,
	}
	if err := uc.repos.User().GetByEmail(ctx, user); err != nil {
		uc.logger.Error("Failed find user by email", zap.Error(err))
		return response.BadRequest(err.Error())
	}

	if err := uc.repos.Auth().VerifyAccount(ctx, otpCode, user); err != nil {
		uc.logger.Error("Failed to use OTP code", zap.Error(err))
		return response.InternalServerError("Failed to verify OTP code")
	}

	return response.Success("Verify Account successfully", 0)
}

func (uc *authUseCase) ResendVerifyAccount(ctx context.Context, payload entity.SendCodeEmail) response.StatusResponse {
	user := &entity.User{
		Email: payload.Email,
	}

	if err := uc.repos.User().GetByEmail(ctx, user); err != nil {
		uc.logger.Error("Failed find user by email", zap.Error(err))
		return response.BadRequest(err.Error())
	}

	if user.IsVerify {
		uc.logger.Error("User verified")
		return response.BadRequest("User verified")
	}

	otpCode := &entity.OtpCode{
		IsUsed:         false,
		OtpCode:        common.GenerateCode(6),
		UserID:         &user.ID,
		IdentifierType: entity.IdentifierTypeEmail.String(),
		Identifier:     payload.Email,
		OtpType:        entity.OtpTypeVerifyAccount.String(),
		ExpiresAt:      common.GetExpireTime(15),
	}

	if err := uc.repos.OtpCode().CreateOtpCode(ctx, otpCode); err != nil {
		uc.logger.Error("Failed create otp code", zap.Error(err))
		return response.InternalServerError("Failed create otp code")
	}

	otpCodeEncrypt, err := helper.Encrypt(otpCode.OtpCode, uc.config.Server.SecretKey)
	if err != nil {
		uc.logger.Error("Failed encrypt otp code", zap.Error(err))
		return response.InternalServerError("Failed encrypt otp code")
	}

	jsonPayload, err := json.Marshal(&entity.SendMailVerifyAccount{
		UserID: user.ID,
		Email:  user.Email,
		Token:  otpCodeEncrypt,
	})
	if err != nil {
		uc.logger.Error("Failed to marshal payload", zap.Error(err))
		return response.InternalServerError("Failed send mail verify account")
	}
	task := asynq.NewTask(string(tasks.TypeSendEmailVerifyAccount), jsonPayload)
	uc.asynqClient.EnqueueContext(ctx, task)

	return response.Success("Resend Verify Account successfully", 0)
}

func (uc *authUseCase) createToken(tokenClaims Claims, accessSecret string) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(accessSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return accessTokenString, nil
}

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

	accessToken, err := uc.createToken(accessTokenClaims, uc.config.JWT.AccessSecret)
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

	refreshToken, err := uc.createToken(refreshTokenClaims, uc.config.JWT.RefreshSecret)
	if err != nil {
		return "", err
	}

	if err := uc.repos.Auth().SetCacheTokenVersion(ctx, entity.RefreshToken, user.ID, string(uuid.String()), uc.config.JWT.RefreshExpiresIn); err != nil {
		uc.logger.Error("Failed to set token version", zap.Error(err))
		return "", fmt.Errorf("failed to set token version: %w", err)
	}

	return refreshToken, nil
}
