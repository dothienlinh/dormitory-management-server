package services

import (
	"dormitory_management/internal/helpers"
	"dormitory_management/internal/models"
	"dormitory_management/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthService struct {
	BaseService
	util *utils.Util
}

func NewAuthService(db *gorm.DB, redis *redis.Client, logger *zap.Logger, util *utils.Util, response *APIResponse) *AuthService {
	return &AuthService{BaseService{db: db, redis: redis, logger: logger, response: response}, util}
}

func (s *AuthService) Login(c *gin.Context, email string, password string) {
	user := models.User{}
	s.db.Where("email = ?", email).First(&user)

	if user.ID == 0 {
		s.logger.Error("User not found", zap.String("email", email))
		s.response.NotFound(c, "User not found")
		return
	}

	if !helpers.VerifyPassword(password, user.Password) {
		s.logger.Error("Invalid password", zap.String("email", email))
		s.response.BadRequest(c, "Invalid password")
		return
	}

	if err := s.util.InvalidateUserTokens(s.redis, user.ID); err != nil {
		s.logger.Error("Failed to invalidate existing tokens", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	accessToken, err := s.util.GenerateAccessToken(user.ID)
	if err != nil {
		s.logger.Error("Error generating access token", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	refreshToken, err := s.util.GenerateRefreshToken(user.ID)
	if err != nil {
		s.logger.Error("Error generating refresh token", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	s.response.Success(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, 0)
}

func (s *AuthService) Register(c *gin.Context, payload models.UserRegister) {
	userExists := models.User{}
	s.db.Where("email = ?", payload.Email).First(&userExists)

	if userExists.ID != 0 {
		s.logger.Error("User already exists", zap.String("email", payload.Email))
		s.response.BadRequest(c, "User already exists")
		return
	}

	user := &models.User{
		FullName: payload.FullName,
		Email:    payload.Email,
		Password: payload.Password,
		Role:     models.UserRoleStudent,
	}

	if err := s.db.Create(&user).Error; err != nil {
		s.logger.Error("Error creating user", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	s.response.Success(c, user, 0)
}

func (s *AuthService) Logout(c *gin.Context, userID uint) {
	if err := s.util.InvalidateUserTokens(s.redis, userID); err != nil {
		s.logger.Error("Failed to invalidate tokens", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	s.response.Success(c, "Logged out successfully", 0)
}

func (s *AuthService) Me(c *gin.Context, userID uint) {
	user := &models.User{}
	if err := s.db.Where("id = ?", userID).First(user).Error; err != nil {
		s.logger.Error("Error getting user", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	s.response.Success(c, user, 0)
}

func (s *AuthService) RefreshToken(c *gin.Context, payload models.UserRefreshToken) {
	claims, err := s.util.ValidateRefreshToken(payload.RefreshToken)
	if err != nil {
		s.logger.Error("Error validating refresh token", zap.Error(err))
		s.response.BadRequest(c, err.Error())
		return
	}

	user := models.User{}
	if err := s.db.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		s.logger.Error("Error getting user", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	if user.ID == 0 {
		s.logger.Error("User not found", zap.Uint("user_id", claims.UserID))
		s.response.NotFound(c, "User not found")
		return
	}

	accessToken, err := s.util.GenerateAccessToken(user.ID)
	if err != nil {
		s.logger.Error("Error generating access token", zap.Error(err))
		s.response.DBError(c, err)
		return
	}

	s.response.Success(c, gin.H{
		"access_token": accessToken,
	}, 0)
}
