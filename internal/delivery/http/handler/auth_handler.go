package handler

import (
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthHandler handles HTTP requests related to authentication
type AuthHandler struct {
	useCases usecase.UseCases
	logger   logger.Logger
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(useCases usecase.UseCases, logger logger.Logger) *AuthHandler {
	return &AuthHandler{
		useCases: useCases,
		logger:   logger,
	}
}

// Register handles the request to register a new user
func (h *AuthHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.logger.Info("Register")

		var user entity.UserRegister
		if err := c.ShouldBindJSON(&user); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
				"error":   "Bad Request",
			})
			return
		}

		resp := h.useCases.Auth().Register(c, &user)
		c.JSON(resp.Status, resp.Response)
	}
}

// Login handles the request to login a user
func (h *AuthHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.logger.Info("Login")

		var loginData entity.UserLogin
		if err := c.ShouldBindJSON(&loginData); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
				"error":   "Bad Request",
			})
			return
		}

		resp := h.useCases.Auth().Login(c, &loginData)
		c.JSON(resp.Status, resp.Response)
	}
}

// RefreshToken handles the request to refresh a token
func (h *AuthHandler) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.logger.Info("RefreshToken")

		var refreshData entity.UserRefreshToken
		if err := c.ShouldBindJSON(&refreshData); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
				"error":   "Bad Request",
			})
			return
		}

		resp := h.useCases.Auth().RefreshToken(c, refreshData.RefreshToken)
		c.JSON(resp.Status, resp.Response)
	}
}

// Logout handles the request to logout a user
func (h *AuthHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.logger.Info("Logout")

		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "User ID not found",
				"error":   "Unauthorized",
			})
			return
		}

		resp := h.useCases.Auth().Logout(c, userID.(uint))
		c.JSON(resp.Status, resp.Response)
	}
}

// Me handles the request to get the current user
func (h *AuthHandler) Me() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.logger.Info("Me")

		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "User ID not found",
				"error":   "Unauthorized",
			})
			return
		}

		resp := h.useCases.Auth().Me(c, userID.(uint))
		c.JSON(resp.Status, resp.Response)
	}
}
