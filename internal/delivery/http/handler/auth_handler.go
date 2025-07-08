package handler

import (
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	useCases usecase.UseCases
	logger   logger.Logger
}

func NewAuthHandler(useCases usecase.UseCases, logger logger.Logger) *AuthHandler {
	return &AuthHandler{
		useCases: useCases,
		logger:   logger,
	}
}

func (h *AuthHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("Register")

		var user entity.UserRegister
		if err := c.ShouldBindJSON(&user); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.Auth().Register(c, &user)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *AuthHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("Login")

		var loginData entity.UserLogin
		if err := c.ShouldBindJSON(&loginData); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.Auth().Login(c, &loginData)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *AuthHandler) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("RefreshToken")

		var refreshData entity.UserRefreshToken
		if err := c.ShouldBindJSON(&refreshData); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.Auth().RefreshToken(c, refreshData.RefreshToken)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *AuthHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("Logout")

		userID, exists := c.Get("userID")
		if !exists {
			resp = response.Unauthorized("User ID not found")
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.Auth().Logout(c, userID.(uint64))
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *AuthHandler) Me() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("Me")

		userID, exists := c.Get("userID")
		if !exists {
			resp = response.Unauthorized("User ID not found")
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.Auth().Me(c, userID.(uint64))
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *AuthHandler) VerifyAccount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("VerifyAccount")

		var payload entity.VerifyAccount
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			ctx.Error(err)
			return
		}

		resp = h.useCases.Auth().VerifyAccount(ctx, payload)
		ctx.JSON(resp.Status, resp.Response)
	}
}

func (h *AuthHandler) ResendVerifyAccount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("ResendVerifyAccount")

		var payload entity.SendCodeEmail
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			ctx.Error(err)
			return
		}

		resp = h.useCases.Auth().ResendVerifyAccount(ctx, payload)
		ctx.JSON(resp.Status, resp.Response)
	}
}

func (h *AuthHandler) ForgotPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("ResendVerifyAccount")

		var payload entity.SendCodeEmail
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			ctx.Error(err)
			return
		}

		resp = h.useCases.Auth().ForgotPassword(ctx, payload)
		ctx.JSON(resp.Status, resp.Response)
	}
}

func (h *AuthHandler) ResetPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("ResetPassword")

		var payload entity.ResetPassword
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			ctx.Error(err)
			return
		}

		resp = h.useCases.Auth().ResetPassword(ctx, payload)
		ctx.JSON(resp.Status, resp.Response)
	}
}

func (h *AuthHandler) ChangePassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("ChangePassword")

		userID, exists := ctx.Get("userID")
		if !exists {
			resp = response.Unauthorized("User ID not found")
			ctx.JSON(resp.Status, resp.Response)
			return
		}

		var payload entity.ChangePassword
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			ctx.Error(err)
			return
		}

		resp = h.useCases.Auth().ChangePassword(ctx, userID.(uint64), payload)
		ctx.JSON(resp.Status, resp.Response)
	}
}
