package handler

import (
	"dormitory_management/internal/domain/response"
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"
	"strconv"

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

		resp := h.usecase.Dashboard().GetStats(ctx)

		ctx.JSON(resp.Status, resp.Response)
	}
}

func (h *DashboardHandler) StudentOverview() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp response.StatusResponse

		userID, exists := ctx.Get("userID")
		if !exists {
			resp = response.Unauthorized("User ID not found")
			ctx.JSON(resp.Status, resp.Response)
			return
		}
		resp = h.usecase.Dashboard().StudentOverview(ctx, userID.(uint64))

		ctx.JSON(resp.Status, resp.Response)
	}
}

// GetNotifications handles getting notifications
func (h *DashboardHandler) GetNotifications() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp response.StatusResponse

		userID, exists := ctx.Get("userID")
		if !exists {
			resp = response.Unauthorized("User ID not found")
			ctx.JSON(resp.Status, resp.Response)
			return
		}

		// Parse query parameters
		page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
		notificationType := ctx.Query("type")

		var isRead *bool
		if isReadStr := ctx.Query("is_read"); isReadStr != "" {
			if isReadValue, err := strconv.ParseBool(isReadStr); err == nil {
				isRead = &isReadValue
			}
		}

		resp = h.usecase.Dashboard().GetNotifications(ctx, userID.(uint64), page, limit, notificationType, isRead)
		ctx.JSON(resp.Status, resp.Response)
	}
}

// MarkNotificationAsRead handles marking a notification as read
func (h *DashboardHandler) MarkNotificationAsRead() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp response.StatusResponse

		userID, exists := ctx.Get("userID")
		if !exists {
			resp = response.Unauthorized("User ID not found")
			ctx.JSON(resp.Status, resp.Response)
			return
		}

		notificationIDStr := ctx.Param("id")
		notificationID, err := strconv.ParseUint(notificationIDStr, 10, 64)
		if err != nil {
			resp = response.BadRequest("Invalid notification ID")
			ctx.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.usecase.Dashboard().MarkNotificationAsRead(ctx, notificationID, userID.(uint64))
		ctx.JSON(resp.Status, resp.Response)
	}
}

// MarkAllNotificationsAsRead handles marking all notifications as read
func (h *DashboardHandler) MarkAllNotificationsAsRead() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp response.StatusResponse

		userID, exists := ctx.Get("userID")
		if !exists {
			resp = response.Unauthorized("User ID not found")
			ctx.JSON(resp.Status, resp.Response)
			return
		}

		var body struct {
			Type string `json:"type"`
		}
		ctx.ShouldBindJSON(&body)

		resp = h.usecase.Dashboard().MarkAllNotificationsAsRead(ctx, userID.(uint64), body.Type)
		ctx.JSON(resp.Status, resp.Response)
	}
}

// GetEvents handles getting events
func (h *DashboardHandler) GetEvents() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Parse query parameters
		page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
		eventType := ctx.Query("type")
		fromDate := ctx.Query("from_date")
		toDate := ctx.Query("to_date")

		resp := h.usecase.Dashboard().GetEvents(ctx, page, limit, eventType, fromDate, toDate)
		ctx.JSON(resp.Status, resp.Response)
	}
}

// GetServiceRequests handles getting service requests
func (h *DashboardHandler) GetServiceRequests() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp response.StatusResponse

		userID, exists := ctx.Get("userID")
		if !exists {
			resp = response.Unauthorized("User ID not found")
			ctx.JSON(resp.Status, resp.Response)
			return
		}

		// Parse query parameters
		page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
		status := ctx.Query("status")
		category := ctx.Query("category")

		resp = h.usecase.Dashboard().GetServiceRequests(ctx, userID.(uint64), page, limit, status, category)
		ctx.JSON(resp.Status, resp.Response)
	}
}

// GetQuickStats handles getting quick stats
func (h *DashboardHandler) GetQuickStats() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp response.StatusResponse

		userID, exists := ctx.Get("userID")
		if !exists {
			resp = response.Unauthorized("User ID not found")
			ctx.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.usecase.Dashboard().GetQuickStats(ctx, userID.(uint64))
		ctx.JSON(resp.Status, resp.Response)
	}
}
