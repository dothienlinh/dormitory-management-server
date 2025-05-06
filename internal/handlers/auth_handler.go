package handlers

import (
	"dormitory_management/internal/helpers"
	"dormitory_management/internal/models"
	"dormitory_management/internal/utils"
	"dormitory_management/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload models.UserRegister
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Register failed", err.Error()))
			return
		}

		if err := pkg.ValidateStruct(payload); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Register failed", err.Error()))
			return
		}

		userExists := models.User{}
		h.dbClient.Where("email = ?", payload.Email).First(&userExists)
		if userExists.ID != 0 {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Register failed", "User already exists"))
			return
		}

		user := models.User{
			FullName: payload.FullName,
			Email:    payload.Email,
			Password: payload.Password,
		}

		h.dbClient.Create(&user)

		c.JSON(http.StatusCreated, utils.SuccessResponse("Register successfully", user))
	}
}

func (h *Handler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload models.UserLogin
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Login failed", err.Error()))
			return
		}

		if err := pkg.ValidateStruct(payload); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Login failed", err.Error()))
			return
		}

		user := models.User{}
		h.dbClient.Where("email = ?", payload.Email).First(&user)
		if user.ID == 0 {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Login failed", "User not found"))
			return
		}

		if !helpers.VerifyPassword(payload.Password, user.Password) {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Login failed", "Invalid password"))
			return
		}

		if err := h.util.InvalidateUserTokens(h.redisClient, user.ID); err != nil {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Login failed", "Failed to invalidate existing tokens"))
			return
		}

		accessToken, err := h.util.GenerateAccessToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Login failed", err.Error()))
			return
		}

		refreshToken, err := h.util.GenerateRefreshToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Login failed", err.Error()))
			return
		}

		c.JSON(http.StatusOK, utils.SuccessResponse("Login successfully", gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		}))
	}
}

func (h *Handler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		if err := h.util.InvalidateUserTokens(h.redisClient, userID); err != nil {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Logout failed", "Failed to invalidate tokens"))
			return
		}

		c.JSON(http.StatusOK, utils.SuccessResponse("Logout successfully", nil))
	}
}

func (h *Handler) Me() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		user := models.User{}
		h.dbClient.Where("id = ?", userID).First(&user)

		c.JSON(http.StatusOK, utils.SuccessResponse("Get current user successfully", user))
	}

}

func (h *Handler) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload models.UserRefreshToken
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Refresh token failed", err.Error()))
			return
		}

		if err := pkg.ValidateStruct(payload); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Refresh token failed", err.Error()))
			return
		}

		claims, err := h.util.ValidateRefreshToken(payload.RefreshToken)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Refresh token failed", err.Error()))
			return
		}

		user := models.User{}
		h.dbClient.Where("id = ?", claims.UserID).First(&user)
		if user.ID == 0 {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Refresh token failed", "User not found"))
			return
		}

		accessToken, err := h.util.GenerateAccessToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Refresh token failed", err.Error()))
			return
		}

		c.JSON(http.StatusOK, utils.SuccessResponse("Refresh token successfully", gin.H{
			"access_token": accessToken,
		}))
	}
}
