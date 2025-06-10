package router

import (
	"dormitory_management/internal/delivery/http/handler"
	"dormitory_management/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes configures user related routes
func SetupUserRoutes(router *gin.RouterGroup, handlers *handler.Handlers, mw *middleware.Middleware) {
	// User routes (requires authentication)
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
