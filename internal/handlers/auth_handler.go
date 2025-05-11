package handlers

import (
	"dormitory_management/internal/models"
	"dormitory_management/pkg"

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

		user, err := h.service.Auth.Register(payload)
		if err != nil {
			h.logger.Error("Error registering user", zap.Error(err))
			BadRequest(c, err.Error())
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

		accessToken, refreshToken, err := h.service.Auth.Login(payload.Email, payload.Password)
		if err != nil {
			h.logger.Error("Error logging in", zap.Error(err))
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

		if err := h.service.Auth.Logout(userID); err != nil {
			h.logger.Error("Error logging out", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		Success(c, "Logged out successfully", 0)
	}
}

func (h *Handler) Me() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		user, err := h.service.Auth.Me(userID)
		if err != nil {
			h.logger.Error("Error getting user", zap.Error(err))
			BadRequest(c, err.Error())
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

		accessToken, err := h.service.Auth.RefreshToken(payload)
		if err != nil {
			h.logger.Error("Error refreshing token", zap.Error(err))
			BadRequest(c, err.Error())
			return
		}

		Success(c, gin.H{
			"access_token": accessToken,
		}, 0)
	}
}
