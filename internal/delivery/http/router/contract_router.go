package router

import (
	"dormitory_management/internal/delivery/http/handler"
	"dormitory_management/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupContractRoutes(router *gin.RouterGroup, handlers *handler.Handlers, mw *middleware.Middleware) {
	contracts := router.Group("/contracts", mw.AuthMiddleware())
	{
		contracts.GET("", handlers.Contract.GetListContracts())
		contracts.GET("/:id", handlers.Contract.GetContractByID())
		contracts.GET("/user/:user_id", handlers.Contract.GetContractByUserID())

		contractAdmin := contracts.Group("", mw.AdminMiddleware())
		{
			contractAdmin.POST("", handlers.Contract.CreateContract())
			contractAdmin.PUT("/:id", handlers.Contract.UpdateContract())
			contractAdmin.DELETE("/:id", handlers.Contract.DeleteContract())
		}
	}
}
