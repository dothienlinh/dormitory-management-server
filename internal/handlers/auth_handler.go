package handlers

import (
	"dormitory_management/internal/database/db"
	"dormitory_management/internal/helpers"
	"dormitory_management/internal/models"
	"dormitory_management/internal/utils"
	"dormitory_management/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
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
	db.DB.Where("email = ?", payload.Email).First(&userExists)
	if userExists.ID != 0 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Register failed", "User already exists"))
		return
	}

	user := models.User{
		FullName: payload.FullName,
		Email:    payload.Email,
		Password: payload.Password,
	}

	db.DB.Create(&user)

	c.JSON(http.StatusCreated, utils.SuccessResponse("Register successfully", user))
}

func Login(c *gin.Context) {
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
	db.DB.Where("email = ?", payload.Email).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Login failed", "User not found"))
		return
	}

	if !helpers.VerifyPassword(payload.Password, user.Password) {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Login failed", "Invalid password"))
		return
	}

	if err := utils.InvalidateUserTokens(user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Login failed", "Failed to invalidate existing tokens"))
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Login failed", err.Error()))
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Login failed", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Login successfully", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}))
}

func Logout(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := utils.InvalidateUserTokens(userID); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Logout failed", "Failed to invalidate tokens"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Logout successfully", nil))
}

func Me(c *gin.Context) {
	userID := c.GetUint("user_id")

	user := models.User{}
	db.DB.Where("id = ?", userID).First(&user)

	c.JSON(http.StatusOK, utils.SuccessResponse("Get current user successfully", user))
}

func RefreshToken(c *gin.Context) {
	var payload models.UserRefreshToken
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Refresh token failed", err.Error()))
		return
	}

	if err := pkg.ValidateStruct(payload); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Refresh token failed", err.Error()))
		return
	}

	claims, err := utils.ValidateRefreshToken(payload.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Refresh token failed", err.Error()))
		return
	}

	user := models.User{}
	db.DB.Where("id = ?", claims.UserID).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Refresh token failed", "User not found"))
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Refresh token failed", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Refresh token successfully", gin.H{
		"access_token": accessToken,
	}))
}
