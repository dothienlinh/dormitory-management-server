package router

import (
	"dormitory_management/internal/delivery/http/handler"
	"dormitory_management/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoomCategoryRoutes(router *gin.RouterGroup, handlers *handler.Handlers, mw *middleware.Middleware) {
	roomCategories := router.Group("/room-categories")
	{
		roomCategories.GET("", handlers.RoomCategory.GetListRoomCategories())
		roomCategories.GET("/:id", handlers.RoomCategory.GetRoomCategoryByID())

		roomCategoryAdmin := roomCategories.Group("", mw.AuthMiddleware(), mw.AdminMiddleware())
		{
			roomCategoryAdmin.POST("", handlers.RoomCategory.CreateRoomCategory())
			roomCategoryAdmin.PUT("/:id", handlers.RoomCategory.UpdateRoomCategory())
			roomCategoryAdmin.DELETE("/:id", handlers.RoomCategory.DeleteRoomCategory())
		}
	}
}
