package server

import (
	"dormitory_management/internal/handlers"
	"dormitory_management/internal/middleware"

	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	handler    *handlers.Handler
	middleware *middleware.Middleware
}

func NewAuthRouter(handler *handlers.Handler, middleware *middleware.Middleware) *AuthRouter {
	return &AuthRouter{handler: handler, middleware: middleware}
}

func (r *AuthRouter) Register(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", r.handler.Login())
		auth.POST("/register", r.handler.Register())
		auth.POST("/refresh-token", r.handler.RefreshToken())
	}
}
