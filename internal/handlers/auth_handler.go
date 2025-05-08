package handlers

import (
	"dormitory_management/internal/helpers"
	"dormitory_management/internal/models"
	"dormitory_management/pkg"
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload models.UserRegister
		if err := c.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Error binding JSON", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		if err := pkg.ValidateStruct(payload); err != nil {
			h.logger.Error("Error validating struct", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		userExists := models.User{}
		if err := h.dbClient.Where("email = ?", payload.Email).First(&userExists).Error; err != nil {
			h.logger.Error("Error checking user existence", zap.Error(err))
			BadRequest(c, errors.New("error checking user existence").Error())
			return
		}

		if userExists.ID != 0 {
			h.logger.Error("User already exists", zap.String("email", payload.Email))
			BadRequest(c, errors.New("user already exists").Error())
			return
		}

		user := models.User{
			FullName: payload.FullName,
			Email:    payload.Email,
			Password: payload.Password,
		}

		if err := h.dbClient.Create(&user).Error; err != nil {
			h.logger.Error("Error creating user", zap.Error(err))
			BadRequest(c, errors.New("error creating user").Error())
			return
		}

		Success(c, user, 0)
	}
}

func (h *Handler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload models.UserLogin
		if err := c.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Error binding JSON", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		if err := pkg.ValidateStruct(payload); err != nil {
			h.logger.Error("Error validating struct", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		user := models.User{}
		if err := h.dbClient.Where("email = ?", payload.Email).First(&user).Error; err != nil {
			h.logger.Error("Error getting user", zap.Error(err))
			BadRequest(c, errors.New("error getting user").Error())
			return
		}

		if user.ID == 0 {
			h.logger.Error("User not found", zap.String("email", payload.Email))
			BadRequest(c, errors.New("user not found").Error())
			return
		}

		if !helpers.VerifyPassword(payload.Password, user.Password) {
			h.logger.Error("Invalid password", zap.String("email", payload.Email))
			BadRequest(c, errors.New("invalid password").Error())
			return
		}

		if err := h.util.InvalidateUserTokens(h.redisClient, user.ID); err != nil {
			h.logger.Error("Failed to invalidate existing tokens", zap.Error(err))
			BadRequest(c, errors.New("failed to invalidate existing tokens").Error())
			return
		}

		accessToken, err := h.util.GenerateAccessToken(user.ID)
		if err != nil {
			h.logger.Error("Error generating access token", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		refreshToken, err := h.util.GenerateRefreshToken(user.ID)
		if err != nil {
			h.logger.Error("Error generating refresh token", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		Success(c, gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		}, 0)
	}
}

func (h *Handler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		if err := h.util.InvalidateUserTokens(h.redisClient, userID); err != nil {
			h.logger.Error("Failed to invalidate tokens", zap.Error(err))
			BadRequest(c, errors.New("failed to invalidate tokens").Error())
			return
		}

		Success(c, nil, 0)
	}
}

func (h *Handler) Me() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		user := models.User{}
		if err := h.dbClient.Where("id = ?", userID).First(&user).Error; err != nil {
			h.logger.Error("Error getting user", zap.Error(err))
			BadRequest(c, errors.New("error getting user").Error())
			return
		}

		Success(c, user, 0)
	}

}

func (h *Handler) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload models.UserRefreshToken
		if err := c.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Error binding JSON", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		if err := pkg.ValidateStruct(payload); err != nil {
			h.logger.Error("Error validating struct", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		claims, err := h.util.ValidateRefreshToken(payload.RefreshToken)
		if err != nil {
			h.logger.Error("Error validating refresh token", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		user := models.User{}
		if err := h.dbClient.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
			h.logger.Error("Error getting user", zap.Error(err))
			BadRequest(c, errors.New("error getting user").Error())
			return
		}

		if user.ID == 0 {
			h.logger.Error("User not found", zap.Uint("user_id", claims.UserID))
			BadRequest(c, errors.New("user not found").Error())
			return
		}

		accessToken, err := h.util.GenerateAccessToken(user.ID)
		if err != nil {
			h.logger.Error("Error generating access token", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		Success(c, gin.H{
			"access_token": accessToken,
		}, 0)
	}
}
