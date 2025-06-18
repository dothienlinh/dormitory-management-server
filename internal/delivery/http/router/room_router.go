package router

import (
	"dormitory_management/internal/delivery/http/handler"
	"dormitory_management/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoomRoutes(router *gin.RouterGroup, handlers *handler.Handlers, mw *middleware.Middleware) {
	rooms := router.Group("/rooms")
	{
		rooms.GET("", handlers.Room.GetListRooms())
		rooms.GET("/:id", handlers.Room.GetRoomByID())

		roomsAdmin := rooms.Group("", mw.AuthMiddleware(), mw.AdminMiddleware())
		{
			roomsAdmin.POST("", handlers.Room.CreateRoom())
			roomsAdmin.PUT("/:id", handlers.Room.UpdateRoom())
			roomsAdmin.DELETE("/:id", handlers.Room.DeleteRoom())
		}
	}
}
