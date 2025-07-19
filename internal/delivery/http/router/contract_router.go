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

		// Student contract routes
		contracts.GET("/my-contract", handlers.Contract.GetMyContract())
		contracts.GET("/:id/download-pdf", handlers.Contract.DownloadContractPDF())
		contracts.GET("/:id/payment-history", handlers.Contract.GetContractPaymentHistory())

		contractAdmin := contracts.Group("", mw.ManagerMiddleware())
		{
			contractAdmin.POST("", handlers.Contract.CreateContract())
			contractAdmin.PUT("/:id", handlers.Contract.UpdateContract())
			contractAdmin.DELETE("/:id", handlers.Contract.DeleteContract())
		}
	}
}
