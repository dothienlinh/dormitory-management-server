package router

import (
	"dormitory_management/internal/delivery/http/handler"
	"dormitory_management/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, handlers *handler.Handlers, mw *middleware.Middleware) {
	apiV1 := router.Group("/api/v1")
	{
		SetupAuthRoutes(apiV1, handlers, mw)
		SetupUserRoutes(apiV1, handlers, mw)
		SetupRoomRoutes(apiV1, handlers, mw)
		SetupRoomCategoryRoutes(apiV1, handlers, mw)
		SetupContractRoutes(apiV1, handlers, mw)
		SetupDashboardRoutes(apiV1, handlers, mw)
		SetupEmailRoutes(apiV1, handlers, mw)
		SetupFacilityRoutes(apiV1, handlers, mw)
	}
}
