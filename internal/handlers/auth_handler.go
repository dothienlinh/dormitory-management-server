package handlers

import (
	"dormitory_management/internal/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload models.UserRegister
		if err := c.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Error binding JSON", zap.Error(err))
			h.response.BadRequest(c, err.Error())
			return
		}

		h.service.Auth.Register(c, payload)
	}
}

func (h *Handler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload models.UserLogin
		if err := c.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Error binding JSON", zap.Error(err))
			h.response.BadRequest(c, err.Error())
			return
		}

		h.service.Auth.Login(c, payload.Email, payload.Password)
	}
}

func (h *Handler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		h.service.Auth.Logout(c, userID)

	}
}

func (h *Handler) Me() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		h.service.Auth.Me(c, userID)

	}

}

func (h *Handler) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload models.UserRefreshToken
		if err := c.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Error binding JSON", zap.Error(err))
			h.response.BadRequest(c, err.Error())
			return
		}

		h.service.Auth.RefreshToken(c, payload)

	}
}
