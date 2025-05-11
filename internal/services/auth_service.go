package services

import (
	"dormitory_management/internal/helpers"
	"dormitory_management/internal/models"
	"dormitory_management/internal/utils"
	"errors"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthService struct {
	db     *gorm.DB
	redis  *redis.Client
	logger *zap.Logger
	util   *utils.Util
}

func NewAuthService(db *gorm.DB, redis *redis.Client, logger *zap.Logger, util *utils.Util) *AuthService {
	return &AuthService{db: db, redis: redis, logger: logger, util: util}
}

func (s *AuthService) Login(email string, password string) (string, string, error) {
	user := models.User{}
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		s.logger.Error("Error getting user", zap.Error(err))
		return "", "", errors.New("error getting user")
	}

	if user.ID == 0 {
		s.logger.Error("User not found", zap.String("email", email))
		return "", "", errors.New("user not found")
	}

	if !helpers.VerifyPassword(password, user.Password) {
		s.logger.Error("Invalid password", zap.String("email", email))
		return "", "", errors.New("invalid password")
	}

	if err := s.util.InvalidateUserTokens(s.redis, user.ID); err != nil {
		s.logger.Error("Failed to invalidate existing tokens", zap.Error(err))
		return "", "", errors.New("failed to invalidate existing tokens")
	}

	accessToken, err := s.util.GenerateAccessToken(user.ID)
	if err != nil {
		s.logger.Error("Error generating access token", zap.Error(err))
		return "", "", errors.New("error generating access token")
	}

	refreshToken, err := s.util.GenerateRefreshToken(user.ID)
	if err != nil {
		s.logger.Error("Error generating refresh token", zap.Error(err))
		return "", "", errors.New("error generating refresh token")
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Register(payload models.UserRegister) (models.User, error) {
	userExists := models.User{}
	s.db.Where("email = ?", payload.Email).First(&userExists)

	if userExists.ID != 0 {
		s.logger.Error("User already exists", zap.String("email", payload.Email))
		return models.User{}, errors.New("user already exists")
	}

	user := models.User{
		FullName: payload.FullName,
		Email:    payload.Email,
		Password: payload.Password,
	}

	if err := s.db.Create(&user).Error; err != nil {
		s.logger.Error("Error creating user", zap.Error(err))
		return models.User{}, errors.New("error creating user")
	}

	return user, nil
}

func (s *AuthService) Logout(userID uint) error {
	if err := s.util.InvalidateUserTokens(s.redis, userID); err != nil {
		s.logger.Error("Failed to invalidate tokens", zap.Error(err))
		return err
	}

	return nil
}

func (s *AuthService) Me(userID uint) (*models.User, error) {
	user := &models.User{}
	if err := s.db.Where("id = ?", userID).First(user).Error; err != nil {
		s.logger.Error("Error getting user", zap.Error(err))
		return nil, errors.New("error getting user")
	}

	return user, nil
}

func (s *AuthService) RefreshToken(payload models.UserRefreshToken) (string, error) {
	claims, err := s.util.ValidateRefreshToken(payload.RefreshToken)
	if err != nil {
		s.logger.Error("Error validating refresh token", zap.Error(err))
		return "", err
	}

	user := models.User{}
	if err := s.db.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		s.logger.Error("Error getting user", zap.Error(err))
		return "", errors.New("error getting user")
	}

	if user.ID == 0 {
		s.logger.Error("User not found", zap.Uint("user_id", claims.UserID))
		return "", errors.New("user not found")
	}

	accessToken, err := s.util.GenerateAccessToken(user.ID)
	if err != nil {
		s.logger.Error("Error generating access token", zap.Error(err))
		return "", err
	}

	return accessToken, nil
}
