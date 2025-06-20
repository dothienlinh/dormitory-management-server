package handler

import (
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MaintenanceHistoryHandler struct {
	useCases usecase.UseCases
	logger   logger.Logger
}

func NewMaintenanceHistoryHandler(useCases usecase.UseCases, logger logger.Logger) *MaintenanceHistoryHandler {
	return &MaintenanceHistoryHandler{
		useCases: useCases,
		logger:   logger,
	}
}

func (h *MaintenanceHistoryHandler) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload entity.CreateMaintenanceHistory

		if err := ctx.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			ctx.Error(err)
			return
		}

		resp := h.useCases.MaintenanceHistory().Create(ctx, &payload)
		ctx.JSON(resp.Status, resp.Response)
	}
}

func (h *MaintenanceHistoryHandler) Detail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp response.StatusResponse
		id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse maintenance history ID", zap.Error(err))
			resp = response.BadRequest("Invalid maintenance history ID")
			ctx.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.MaintenanceHistory().Detail(ctx, id)
		ctx.JSON(resp.Status, resp.Response)
	}
}

func (h *MaintenanceHistoryHandler) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp response.StatusResponse
		id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse maintenance history ID", zap.Error(err))
			resp = response.BadRequest("Invalid maintenance history ID")
			ctx.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.MaintenanceHistory().Update(ctx, id)
		ctx.JSON(resp.Status, resp.Response)
	}
}

func (h *MaintenanceHistoryHandler) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp response.StatusResponse
		id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse maintenance history ID", zap.Error(err))
			resp = response.BadRequest("Invalid maintenance history ID")
			ctx.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.MaintenanceHistory().Delete(ctx, id)
		ctx.JSON(resp.Status, resp.Response)
	}
}
