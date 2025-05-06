package middleware

import (
	"dormitory_management/internal/models"
	"dormitory_management/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) RoleAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetUint("user_id")

		user := models.User{}
		m.dbClient.Where("id = ?", userId).First(&user)

		if user.Role != models.UserRoleAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, utils.ErrorResponse("User is not an admin", nil))
			return
		}

		c.Next()
	}

}
