package router

import (
	"dormitory_management/internal/delivery/http/handler"
	"dormitory_management/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.RouterGroup, handlers *handler.Handlers, mw *middleware.Middleware) {
	users := router.Group("/users", mw.AuthMiddleware())
	{
		users.GET("", handlers.User.GetListUsers())
		users.GET("/:id", handlers.User.GetUserByID())
		users.PUT("/:id", handlers.User.UpdateUser())
		users.DELETE("/:id", mw.AdminMiddleware(), handlers.User.DeleteUser())
		users.POST("/room", handlers.User.AddUserToRoom())
		users.POST("/remove-room", handlers.User.UserLeavesRoom())
	}
}
