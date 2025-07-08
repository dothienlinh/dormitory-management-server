package router

import (
	"dormitory_management/internal/delivery/http/handler"
	"dormitory_management/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.RouterGroup, handlers *handler.Handlers, mw *middleware.Middleware) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", handlers.Auth.Register())
		auth.POST("/login", handlers.Auth.Login())
		auth.POST("/refresh-token", handlers.Auth.RefreshToken())
		auth.POST("/logout", mw.AuthMiddleware(), handlers.Auth.Logout())
		auth.GET("/me", mw.AuthMiddleware(), handlers.Auth.Me())
		auth.POST("/verify-account", handlers.Auth.VerifyAccount())
		auth.POST("/resend-verify-account", handlers.Auth.ResendVerifyAccount())
		auth.POST("/forgot-password", handlers.Auth.ForgotPassword())
		auth.POST("/reset-password", handlers.Auth.ResetPassword())
		auth.POST("/change-password", mw.AuthMiddleware(), handlers.Auth.ChangePassword())
	}
}
