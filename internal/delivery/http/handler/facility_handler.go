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

type FacilityHandler struct {
	useCases usecase.UseCases
	logger   logger.Logger
}

func NewFacilityHandler(useCases usecase.UseCases, logger logger.Logger) *FacilityHandler {
	return &FacilityHandler{
		useCases: useCases,
		logger:   logger,
	}
}

func (h *FacilityHandler) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp response.StatusResponse

		var payload entity.CreateFacility
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			ctx.Error(err)
			return
		}

		resp = h.useCases.Facilities().Create(ctx, &payload)

		ctx.JSON(resp.Status, resp.Response)
	}
}

func (h *FacilityHandler) List() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp := h.useCases.Facilities().List(ctx)

		ctx.JSON(resp.Status, resp.Response)
	}
}

func (h *FacilityHandler) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp response.StatusResponse

		id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse room category ID", zap.Error(err))
			resp = response.BadRequest("Invalid room category ID")
			ctx.JSON(resp.Status, resp.Response)
			return
		}

		var payload entity.UpdateFacility
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			ctx.Error(err)
			return
		}

		resp = h.useCases.Facilities().Update(ctx, &payload, id)

		ctx.JSON(resp.Status, resp.Response)
	}
}

func (h *FacilityHandler) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp response.StatusResponse

		id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse room category ID", zap.Error(err))
			resp = response.BadRequest("Invalid room category ID")
			ctx.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.Facilities().Delete(ctx, id)

		ctx.JSON(resp.Status, resp.Response)
	}
}
