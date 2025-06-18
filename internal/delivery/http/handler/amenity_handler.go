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

type AmenityHandler struct {
	useCases usecase.UseCases
	logger   logger.Logger
}

func NewAmenityHandler(useCases usecase.UseCases, logger logger.Logger) *AmenityHandler {
	return &AmenityHandler{
		useCases: useCases,
		logger:   logger,
	}
}

func (h *AmenityHandler) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp response.StatusResponse

		var payload entity.CreateAmenity
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			ctx.Error(err)
			return
		}

		resp = h.useCases.Amenities().Create(ctx, &payload)

		ctx.JSON(resp.Status, resp.Response)
	}
}

func (h *AmenityHandler) List() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp := h.useCases.Amenities().List(ctx)

		ctx.JSON(resp.Status, resp.Response)
	}
}

func (h *AmenityHandler) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp response.StatusResponse

		id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse room category ID", zap.Error(err))
			resp = response.BadRequest("Invalid room category ID")
			ctx.JSON(resp.Status, resp.Response)
			return
		}

		var payload entity.UpdateAmenity
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			ctx.Error(err)
			return
		}

		resp = h.useCases.Amenities().Update(ctx, &payload, id)

		ctx.JSON(resp.Status, resp.Response)
	}
}

func (h *AmenityHandler) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp response.StatusResponse

		id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse room category ID", zap.Error(err))
			resp = response.BadRequest("Invalid room category ID")
			ctx.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.Amenities().Delete(ctx, id)

		ctx.JSON(resp.Status, resp.Response)
	}
}
