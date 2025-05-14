package router

import (
	"dormitory_management/internal/delivery/http/handler"
	"dormitory_management/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoomCategoryRoutes configures room category related routes
func SetupRoomCategoryRoutes(router *gin.RouterGroup, handlers *handler.Handlers, mw *middleware.Middleware) {
	// Room category routes
	roomCategories := router.Group("/room-categories")
	{
		roomCategories.GET("", handlers.RoomCategory.GetListRoomCategories())
		roomCategories.GET("/:id", handlers.RoomCategory.GetRoomCategoryByID())

		// These routes require authentication and admin role
		roomCategoryAdmin := roomCategories.Group("", mw.AuthMiddleware(), mw.AdminMiddleware())
		{
			roomCategoryAdmin.POST("", handlers.RoomCategory.CreateRoomCategory())
			roomCategoryAdmin.PUT("/:id", handlers.RoomCategory.UpdateRoomCategory())
			roomCategoryAdmin.DELETE("/:id", handlers.RoomCategory.DeleteRoomCategory())
		}
	}
}
