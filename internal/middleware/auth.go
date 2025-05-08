package middleware

import (
	"dormitory_management/internal/handlers"
	"dormitory_management/internal/models"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			m.logger.Error("Authorization header is required")
			handlers.BadRequest(c, errors.New("authorization header is required").Error())
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := m.util.ValidateAccessToken(token)
		if err != nil {
			m.logger.Error("Error validating access token", zap.Error(err))
			handlers.BadRequest(c, err.Error())
			c.Abort()
			return
		}

		user := models.User{}
		if err := m.dbClient.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
			m.logger.Error("Error getting user", zap.Error(err))
			handlers.BadRequest(c, errors.New("error getting user").Error())
			c.Abort()
			return
		}

		if user.ID == 0 {
			m.logger.Error("User not found", zap.Uint("user_id", claims.UserID))
			handlers.BadRequest(c, errors.New("user not found").Error())
			c.Abort()
			return
		}

		if user.Status != models.UserStatusActive {
			m.logger.Error("User is not active", zap.Uint("user_id", claims.UserID))
			handlers.BadRequest(c, errors.New("user is not active").Error())
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
