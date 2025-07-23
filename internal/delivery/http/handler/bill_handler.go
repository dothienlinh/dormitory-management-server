package handler

import (
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BillHandler struct {
	useCases usecase.UseCases
	logger   logger.Logger
}

func NewBillHandler(useCases usecase.UseCases, logger logger.Logger) *BillHandler {
	return &BillHandler{
		useCases: useCases,
		logger:   logger,
	}
}

func (h *BillHandler) MyListBills() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp response.StatusResponse

		userID, exists := ctx.Get("user_id")
		if !exists {
			resp = response.Unauthorized("User ID not found")
			ctx.JSON(resp.Status, resp.Response)
			return
		}

		var query entity.QueryBill
		if err := ctx.ShouldBindQuery(&query); err != nil {
			h.logger.Error("Failed to bind query parameters", zap.Error(err))
			ctx.Error(err)
			return
		}

		bills, total, err := h.useCases.Bill().MyListBills(ctx, userID.(uint64), &query)
		if err != nil {
			h.logger.Error("Failed to retrieve bills", zap.Error(err))
			resp = response.InternalServerError("Failed to retrieve bills")
			ctx.JSON(resp.Status, resp.Response)
			return
		}

		resp = response.Success(bills, total)
		ctx.JSON(resp.Status, resp.Response)
	}
}
