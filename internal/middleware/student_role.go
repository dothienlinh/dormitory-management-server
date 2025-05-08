package middleware

import (
	"dormitory_management/internal/handlers"
	"dormitory_management/internal/models"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) RoleStudentMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetUint("user_id")

		user := models.User{}
		m.dbClient.Where("id = ?", userId).First(&user)

		if user.Role != models.UserRoleStudent {
			handlers.Forbidden(c, "Forbidden")
			c.Abort()
			return
		}

		c.Next()
	}
}
