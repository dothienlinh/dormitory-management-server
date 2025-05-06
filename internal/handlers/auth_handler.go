package handlers

import (
	"dormitory_management/internal/helpers"
	"dormitory_management/internal/models"
	"dormitory_management/pkg"
	"errors"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload models.UserRegister
		if err := c.ShouldBindJSON(&payload); err != nil {
			BadRequest(c, err.Error())
			return
		}

		if err := pkg.ValidateStruct(payload); err != nil {
			BadRequest(c, err.Error())
			return
		}

		userExists := models.User{}
		h.dbClient.Where("email = ?", payload.Email).First(&userExists)
		if userExists.ID != 0 {
			BadRequest(c, errors.New("user already exists").Error())
			return
		}

		user := models.User{
			FullName: payload.FullName,
			Email:    payload.Email,
			Password: payload.Password,
		}

		h.dbClient.Create(&user)

		Success(c, user, 0)
	}
}

func (h *Handler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload models.UserLogin
		if err := c.ShouldBindJSON(&payload); err != nil {
			BadRequest(c, err.Error())
			return
		}

		if err := pkg.ValidateStruct(payload); err != nil {
			BadRequest(c, err.Error())
			return
		}

		user := models.User{}
		h.dbClient.Where("email = ?", payload.Email).First(&user)
		if user.ID == 0 {
			BadRequest(c, errors.New("user not found").Error())
			return
		}

		if !helpers.VerifyPassword(payload.Password, user.Password) {
			BadRequest(c, errors.New("invalid password").Error())
			return
		}

		if err := h.util.InvalidateUserTokens(h.redisClient, user.ID); err != nil {
			BadRequest(c, errors.New("failed to invalidate existing tokens").Error())
			return
		}

		accessToken, err := h.util.GenerateAccessToken(user.ID)
		if err != nil {
			BadRequest(c, err.Error())
			return
		}

		refreshToken, err := h.util.GenerateRefreshToken(user.ID)
		if err != nil {
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
		h.dbClient.Where("id = ?", userID).First(&user)

		Success(c, user, 0)
	}

}

func (h *Handler) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload models.UserRefreshToken
		if err := c.ShouldBindJSON(&payload); err != nil {
			BadRequest(c, err.Error())
			return
		}

		if err := pkg.ValidateStruct(payload); err != nil {
			BadRequest(c, err.Error())
			return
		}

		claims, err := h.util.ValidateRefreshToken(payload.RefreshToken)
		if err != nil {
			BadRequest(c, err.Error())
			return
		}

		user := models.User{}
		h.dbClient.Where("id = ?", claims.UserID).First(&user)
		if user.ID == 0 {
			BadRequest(c, errors.New("user not found").Error())
			return
		}

		accessToken, err := h.util.GenerateAccessToken(user.ID)
		if err != nil {
			BadRequest(c, err.Error())
			return
		}

		Success(c, gin.H{
			"access_token": accessToken,
		}, 0)
	}
}
