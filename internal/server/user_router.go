package server

import (
	"dormitory_management/internal/handlers"
	"dormitory_management/internal/middleware"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	handler    *handlers.Handler
	middleware *middleware.Middleware
}

func NewUserRouter(handler *handlers.Handler, middleware *middleware.Middleware) *UserRouter {
	return &UserRouter{handler: handler, middleware: middleware}
}

func (r *UserRouter) Register(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.GET("/me", r.handler.Me())
		auth.POST("/logout", r.handler.Logout())
	}
}
