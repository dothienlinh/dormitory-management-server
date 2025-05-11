package server

import (
	"dormitory_management/internal/handlers"
	"dormitory_management/internal/middleware"

	"github.com/gin-gonic/gin"
)

type AdminRouter struct {
	handler    *handlers.Handler
	middleware *middleware.Middleware
}

func NewAdminRouter(handler *handlers.Handler, middleware *middleware.Middleware) *AdminRouter {
	return &AdminRouter{handler: handler, middleware: middleware}
}

func (r *AdminRouter) Register(router *gin.RouterGroup) {
	admin := router.Group("/admin")
	admin.Use(r.middleware.RoleAdminMiddleware())
	{
		rooms := admin.Group("/rooms")
		{
			rooms.POST("/", r.handler.CreateRoom())
			rooms.PUT("/:id", r.handler.ValidateRoom(), r.handler.UpdateRoom())
			rooms.DELETE("/:id", r.handler.ValidateRoom(), r.handler.DeleteRoom())
			rooms.GET("/:id/students", r.handler.GetListStudentsInRoom())
		}

		roomCategories := admin.Group("/room-categories")
		{
			roomCategories.POST("/", r.handler.CreateRoomCategory())
			roomCategories.PUT("/:id", r.handler.ValidateRoomCategory(), r.handler.UpdateRoomCategory())
			roomCategories.DELETE("/:id", r.handler.ValidateRoomCategory(), r.handler.DeleteRoomCategory())
		}

		users := admin.Group("/users")
		{
			users.GET("/", r.handler.GetListUser())
			users.POST("/:user_id/add-to-room/:room_id", r.handler.AddUserToRoom())
			users.DELETE("/:user_id/remove-from-room/:room_id", r.handler.RemoveUserFromRoom())
		}

		contract := admin.Group("/contract")
		{
			contract.POST("/", r.handler.CreateContract())
		}
	}
}
