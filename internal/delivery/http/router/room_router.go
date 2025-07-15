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

		roomsAdmin := rooms.Group("", mw.AuthMiddleware(), mw.ManagerMiddleware())
		{
			roomsAdmin.POST("", handlers.Room.CreateRoom())
			roomsAdmin.PUT("/:id", handlers.Room.UpdateRoom())
			roomsAdmin.DELETE("/:id", handlers.Room.DeleteRoom())
		}
	}

	// Student Room Routes
	studentRoom := router.Group("/student/room", mw.AuthMiddleware())
	{
		studentRoom.GET("/details", handlers.Room.GetStudentRoomDetails())
		studentRoom.GET("/stats", handlers.Room.GetRoomStats())
		studentRoom.GET("/roommates", handlers.Room.GetRoommates())

		// Room Issues
		studentRoom.GET("/issues", handlers.Room.GetRoomIssues())
		studentRoom.POST("/issues", handlers.Room.CreateRoomIssue())
		studentRoom.GET("/issues/:issueId", handlers.Room.GetRoomIssueDetails())

		// Room Bills
		studentRoom.GET("/bills", handlers.Room.GetRoomBills())

		// Room Rules and Cleaning
		studentRoom.GET("/rules", handlers.Room.GetRoomRules())
		studentRoom.GET("/cleaning-schedule", handlers.Room.GetCleaningSchedule())
	}
}
