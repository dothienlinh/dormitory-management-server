package middleware

import (
	"dormitory_management/internal/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (m *Middleware) RoleAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetUint("user_id")

		user := models.User{}
		if err := m.dbClient.Where("id = ?", userId).First(&user).Error; err != nil {
			m.logger.Error("Error getting user", zap.Error(err))
			m.response.Forbidden(c, "Forbidden")
			c.Abort()
			return
		}

		if user.Role != models.UserRoleAdmin {
			m.logger.Error("User is not admin", zap.Uint("user_id", userId))
			m.response.Forbidden(c, "Forbidden")
			c.Abort()
			return
		}

		c.Next()
	}

}
