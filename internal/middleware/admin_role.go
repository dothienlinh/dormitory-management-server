package middleware

import (
	"dormitory_management/internal/handlers"
	"dormitory_management/internal/models"
	"errors"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) RoleAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetUint("user_id")

		user := models.User{}
		m.dbClient.Where("id = ?", userId).First(&user)

		if user.Role != models.UserRoleAdmin {
			handlers.BadRequest(c, errors.New("user is not an admin").Error())
			return
		}

		c.Next()
	}

}
