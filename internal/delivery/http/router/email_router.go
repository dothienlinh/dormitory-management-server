package router

import (
	"dormitory_management/internal/delivery/http/handler"
	"dormitory_management/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupEmailRoutes(router *gin.RouterGroup, handlers *handler.Handlers, mw *middleware.Middleware) {
	email := router.Group("/email")
	{
		email.POST("/send-code", handlers.Email.SendEmail())
		email.POST("/verify-code", handlers.Email.VerifyCode())
	}
}
