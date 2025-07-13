package router

import (
	"dormitory_management/internal/delivery/http/handler"
	"dormitory_management/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.RouterGroup, handlers *handler.Handlers, mw *middleware.Middleware) {
	users := router.Group("/users", mw.AuthMiddleware())
	{
		users.GET("", mw.ManagerMiddleware(), handlers.User.GetListUsers())
		users.GET("/:id", mw.ManagerMiddleware(), handlers.User.GetUserByID())
		users.PUT("/me", mw.AuthMiddleware(), handlers.User.UpdateMe())
		users.PUT("/:id", mw.ManagerMiddleware(), handlers.User.UpdateUser())
		users.DELETE("/:id", mw.AdminMiddleware(), handlers.User.DeleteUser())
		users.POST("/room", mw.StaffMiddleware(), handlers.User.AddUserToRoom())
		users.POST("/remove-room", mw.StaffMiddleware(), handlers.User.UserLeavesRoom())
		users.PUT("/:id/status-account", mw.AdminMiddleware(), handlers.User.UpdateUserStatusAccount())
	}
}
