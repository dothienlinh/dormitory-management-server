package handler

import (
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	usecase usecase.UseCases
	logger  logger.Logger
}

func NewDashboardHandler(usecase usecase.UseCases, logger logger.Logger) *DashboardHandler {
	return &DashboardHandler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *DashboardHandler) GetStats() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h.logger.Info("GetStats")

		resp := h.usecase.Dashboard().GetStats(ctx)

		ctx.JSON(resp.Status, resp.Response)
	}
}
