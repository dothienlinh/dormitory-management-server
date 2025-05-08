package server

import (
	"dormitory_management/internal/handlers"
	"dormitory_management/internal/middleware"

	"github.com/gin-gonic/gin"
)

type StudentRouter struct {
	handler    *handlers.Handler
	middleware *middleware.Middleware
}

func NewStudentRouter(handler *handlers.Handler, middleware *middleware.Middleware) *StudentRouter {
	return &StudentRouter{handler: handler, middleware: middleware}
}

func (r *StudentRouter) Register(router *gin.RouterGroup) {
	student := router.Group("/student")
	student.Use(r.middleware.RoleStudentMiddleware())
	{
		//TODO: Implement student router
	}
}
