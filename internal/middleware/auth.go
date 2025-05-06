package middleware

import (
	"dormitory_management/internal/handlers"
	"dormitory_management/internal/models"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			handlers.BadRequest(c, errors.New("authorization header is required").Error())
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := m.util.ValidateAccessToken(token)
		if err != nil {
			handlers.BadRequest(c, err.Error())
			c.Abort()
			return
		}

		user := models.User{}
		m.dbClient.Where("id = ?", claims.UserID).First(&user)

		if user.ID == 0 {
			handlers.BadRequest(c, errors.New("user not found").Error())
			c.Abort()
			return
		}

		if user.Status != models.UserStatusActive {
			handlers.BadRequest(c, errors.New("user is not active").Error())
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
