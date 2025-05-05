package middleware

import (
	"dormitory_management/internal/database/db"
	"dormitory_management/internal/models"
	"dormitory_management/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Authorization header is required", nil))
		c.Abort()
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := utils.ValidateAccessToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid token", err.Error()))
		c.Abort()
		return
	}

	user := models.User{}
	db.DB.Where("id = ?", claims.UserID).First(&user)

	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not found", nil))
		c.Abort()
		return
	}

	if user.Status != models.UserStatusActive {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User is not active", nil))
		c.Abort()
		return
	}

	c.Set("user_id", claims.UserID)
	c.Next()
}
