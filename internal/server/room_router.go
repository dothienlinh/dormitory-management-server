package server

import (
	"dormitory_management/internal/handlers"
	"dormitory_management/internal/middleware"

	"github.com/gin-gonic/gin"
)

type RoomRouter struct {
	handler    *handlers.Handler
	middleware *middleware.Middleware
}

func NewRoomRouter(handler *handlers.Handler, middleware *middleware.Middleware) *RoomRouter {
	return &RoomRouter{handler: handler, middleware: middleware}
}

func (r *RoomRouter) Register(router *gin.RouterGroup) {
	router.GET("/rooms", r.handler.GetRooms())
	router.GET("/rooms/:id", r.handler.GetRoomDetail())

	router.GET("/room-categories", r.handler.GetRoomCategories())
	router.GET("/room-categories/:id", r.handler.GetRoomCategoryDetail())
}
