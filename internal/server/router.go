package server

import (
	"dormitory_management/internal/handlers"
	"dormitory_management/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Router interface {
	Register(router *gin.RouterGroup)
}

func NewRouter(middleware *middleware.Middleware, handler *handlers.Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(handlers.ErrorHandler())

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{

			authRouter := NewAuthRouter(handler, middleware)
			authRouter.Register(v1)

			protected := v1.Group("/")
			protected.Use(middleware.AuthMiddleware())
			{
				userRouter := NewUserRouter(handler, middleware)
				userRouter.Register(protected)

				roomRouter := NewRoomRouter(handler, middleware)
				roomRouter.Register(protected)

				adminRouter := NewAdminRouter(handler, middleware)
				adminRouter.Register(protected)

				studentRouter := NewStudentRouter(handler, middleware)
				studentRouter.Register(protected)
			}
		}
	}

	return router
}
